package strength

import (
	"testing"
)

func Test_calcCubicEquation(t *testing.T) {
	type args struct {
		a float64
		b float64
		c float64
		d float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		want1   complex128
		want2   complex128
		wantErr bool
	}{
		{
			name:    "Q >= 0",
			args:    args{a: 1, b: 1, c: 1, d: 1},
			want:    -1,
			want1:   5.551115123125783e-17 + 0.9999999999999999i,
			want2:   5.551115123125783e-17 - 0.9999999999999999i,
			wantErr: false,
		},
		{
			name:    "Q < 0",
			args:    args{a: 1, b: -2, c: -1, d: 1},
			want:    2.246979603717467,
			want1:   0.5549581320873711 + 0i,
			want2:   -0.8019377358048384 + 0i,
			wantErr: false,
		},
		{
			name:    "not cubic equation",
			args:    args{a: 0, b: -2, c: -1, d: 1},
			want:    0,
			want1:   0 + 0i,
			want2:   0 + 0i,
			wantErr: true,
		},
		{
			name:    "test exampl",
			args:    args{a: 1, b: -0.24, c: 0, d: -1.57},
			want:    1.2480088297460978,
			want1:   -0.5040044148730489 + 1.0019897553251107i,
			want2:   -0.5040044148730489 - 1.0019897553251107i,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := calcCubicEquation(tt.args.a, tt.args.b, tt.args.c, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("calcCubicEquation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calcCubicEquation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("calcCubicEquation() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("calcCubicEquation() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_calcKappa(t *testing.T) {
	type args struct {
		relations float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "minimum test",
			args:    args{0.5},
			wantErr: true,
		},
		{
			name: "low test",
			args: args{1},
			want: 0.0138,
		},
		{
			name: "maximum test",
			args: args{12},
			want: 0.0284,
		},
		{
			name: "medium test",
			args: args{2.5},
			want: 0.02775,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calcKappa(tt.args.relations)
			if (err != nil) != tt.wantErr {
				t.Errorf("calcKappa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calcKappa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcRho(t *testing.T) {
	tests := []struct {
		name string
		a    Approx
		want float64
	}{
		{
			name: "free plate press is zero",
			a: Approx{
				width:    120,
				length:   60,
				pressure: 0,
			},
			want: 1,
		},
		{
			name: "free plate press not zero",
			a: Approx{
				width:    120,
				length:   60,
				pressure: 1,
			},
			want: 2.9299999999999997,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcRho(); got != tt.want {
				t.Errorf("Approx.calcRho() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcStartCurvature(t *testing.T) {
	tests := []struct {
		name string
		a    Approx
		want float64
	}{
		{
			name: "test 1",
			a: Approx{
				length:    60,
				thickness: 10,
			},
			want: 0.55,
		},
		{
			name: "test 2",
			a: Approx{
				length:    85.9,
				thickness: 0.6,
			},
			want: 4.151833333333334,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcStartCurvature(); got != tt.want {
				t.Errorf("Approx.calcStartCurvature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcEulerianStrain(t *testing.T) {
	tests := []struct {
		name string
		a    Approx
		want float64
	}{

		{
			name: "length > width",
			a: Approx{
				width:     60,
				length:    120,
				thickness: 0.6,
			},
			want: 0.07600000000000001,
		},
		{
			name: "length < width",
			a: Approx{
				width:     120,
				length:    60,
				thickness: 0.9,
			},
			want: 0.06679687499999999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcEulerianStrain(); got != tt.want {
				t.Errorf("Approx.calcEulerianStrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcPressCurvature(t *testing.T) {
	type args struct {
		k float64
		e float64
	}
	kappa, _ := calcKappa(120 / 60)
	tests := []struct {
		name string
		a    Approx
		args args
		want float64
	}{
		{
			name: "curvature test",
			a: Approx{
				width:     120,
				length:    60,
				thickness: 0.6,
				pressure:  60,
			},
			args: args{
				k: kappa,
				e: 200000,
			},
			want: 496800.00000000006,
		},
		{
			name: "curvature zero press test",
			a: Approx{
				width:     120,
				length:    60,
				thickness: 0.6,
				pressure:  0,
			},
			args: args{
				k: kappa,
				e: 200000,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcPressCurvature(tt.args.k, tt.args.e); got != tt.want {
				t.Errorf("Approx.calcPressCurvature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcChainStrain(t *testing.T) {
	type args struct {
		x              float64
		rho            float64
		eulerianStrain float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test chain strain 1",
			args: args{0, 1, 120},
			want: -120,
		},
		{
			name: "test chain strain 2",
			args: args{1, 1, 120},
			want: 0,
		},
		{
			name: "test chain strain 2",
			args: args{2, 1, 120},
			want: 120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcChainStrain(tt.args.x, tt.args.rho, tt.args.eulerianStrain); got != tt.want {
				t.Errorf("calcChainStrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_momentOfResistance(t *testing.T) {
	type args struct {
		momentOfInertia float64
		centerOfMass    float64
		height          float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "centerOfMass > height",
			args: args{6000, 4, 2},
			want: 3000,
		},
		{
			name: "centerOfMass < height",
			args: args{6000, 4, 8},
			want: 1500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := momentOfResistance(tt.args.momentOfInertia, tt.args.centerOfMass, tt.args.height); got != tt.want {
				t.Errorf("momentOfResistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcActualStrain(t *testing.T) {
	type args struct {
		height          float64
		moment          float64
		centerOfMass    float64
		momentOfInertia float64
		momentFlag      bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "прогиб в сжатой зоне",
			args: args{
				height:          3,
				centerOfMass:    1.5,
				momentFlag:      false,
				moment:          3000,
				momentOfInertia: 7000,
			},
			want: -0.6428571428571428,
		},
		{
			name: "прогиб в растянутой зоне",
			args: args{
				height:          0.75,
				centerOfMass:    1.5,
				momentFlag:      false,
				moment:          3000,
				momentOfInertia: 7000,
			},

			want: 0.3214285714285714,
		},
		{
			name: "перегиб в сжатой зоне",
			args: args{
				height:          0.75,
				centerOfMass:    1.5,
				momentFlag:      true,
				moment:          3000,
				momentOfInertia: 7000,
			},
			want: -0.3214285714285714,
		},
		{
			name: "перегиб в растянутой зоне",
			args: args{
				height:          3,
				centerOfMass:    1.5,
				momentFlag:      true,
				moment:          3000,
				momentOfInertia: 7000,
			},
			want: 0.6428571428571428,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcActualStrain(tt.args.height, tt.args.moment, tt.args.centerOfMass, tt.args.momentOfInertia, tt.args.momentFlag); got != tt.want {
				t.Errorf("calcActualStrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcX(t *testing.T) {
	type args struct {
		rho       float64
		startCurv float64
		pressCurv float64
		eulStrain float64
		actStrain float64
	}
	tests := []struct {
		name string
		a    Approx
		args args
		want float64
	}{
		{
			name: "first case",
			a: Approx{
				length:    60,
				width:     120,
				thickness: 5,
			},
			args: args{
				startCurv: 0.7,
				pressCurv: 0.086,
				eulStrain: 2.16,
				rho:       2.93,
				actStrain: 3.2,
			},
			want: 1.2620423763241606,
		},
		{
			name: "second case",
			a: Approx{
				length:    60,
				width:     120,
				thickness: 6,
			},
			args: args{
				startCurv: 0.6,
				pressCurv: 0,
				eulStrain: 312,
				rho:       1,
				actStrain: -612,
			},
			want: 0.7144262632835598,
		},
		{
			name: "tred case",
			a: Approx{
				length:    60,
				width:     100,
				thickness: 5.5,
			},
			args: args{
				startCurv: 0.6,
				pressCurv: 0.09,
				eulStrain: 310,
				rho:       2.79,
				actStrain: 425,
			},
			want: 1.3314665137022925,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcX(tt.args.rho, tt.args.startCurv, tt.args.pressCurv, tt.args.eulStrain, tt.args.actStrain); got != tt.want {
				t.Errorf("Approx.calcX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApprox_calcReducing(t *testing.T) {
	type args struct {
		actStrain    float64
		elasticModul float64
		startCurv    float64
	}
	tests := []struct {
		name string
		a    Approx
		args args
		want float64
	}{
		{
			name: "first case",
			a: Approx{
				length:    60,
				width:     59,
				thickness: 6,
			},
			args: args{
				actStrain: -10.59118,
			},
			want: 0.742109037832474,
		},
		{
			name: "second case",
			a: Approx{
				length:    60,
				width:     120,
				thickness: 6,
			},
			args: args{
				actStrain: -6.00167,
				startCurv: 0.6,
			},
			want: 0.14403822881165262,
		},
		{
			name: "tred case",
			a: Approx{
				length:    60,
				width:     100,
				thickness: 5.5,
				pressure:  0.09,
			},
			args: args{
				actStrain:    4.16783,
				startCurv:    0.6,
				elasticModul: 2000000,
			},
			want: 0.6728590897093516,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.calcReducing(tt.args.actStrain, tt.args.startCurv, tt.args.elasticModul); got != tt.want {
				t.Errorf("Approx.calcReducing() = %v, want %v", got, tt.want)
			}
		})
	}
}
