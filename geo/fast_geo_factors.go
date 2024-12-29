package geo

import "math"

// Pre-computed meters per degree of latitude, 0 to 90 degrees.
var metersPerDegreeLat = []float64{
	110574.272700,
	110574.610878,
	110575.625009,
	110577.313889,
	110579.675511,
	110582.707071,
	110586.404965,
	110590.764801,
	110595.781395,
	110601.448784,
	110607.760229,
	110614.708223,
	110622.284500,
	110630.480041,
	110639.285090,
	110648.689158,
	110658.681041,
	110669.248827,
	110680.379912,
	110692.061014,
	110704.278186,
	110717.016837,
	110730.261739,
	110743.997053,
	110758.206344,
	110772.872596,
	110787.978236,
	110803.505153,
	110819.434716,
	110835.747799,
	110852.424800,
	110869.445665,
	110886.789912,
	110904.436651,
	110922.364614,
	110940.552173,
	110958.977372,
	110977.617948,
	110996.451360,
	111015.454813,
	111034.605288,
	111053.879568,
	111073.254263,
	111092.705844,
	111112.210666,
	111131.745000,
	111151.285058,
	111170.807026,
	111190.287090,
	111209.701467,
	111229.026434,
	111248.238355,
	111267.313713,
	111286.229139,
	111304.961438,
	111323.487623,
	111341.784938,
	111359.830892,
	111377.603283,
	111395.080231,
	111412.240200,
	111429.062029,
	111445.524959,
	111461.608658,
	111477.293248,
	111492.559331,
	111507.388014,
	111521.760933,
	111535.660275,
	111549.068805,
	111561.969887,
	111574.347503,
	111586.186278,
	111597.471498,
	111608.189131,
	111618.325842,
	111627.869014,
	111636.806763,
	111645.127957,
	111652.822225,
	111659.879975,
	111666.292406,
	111672.051518,
	111677.150126,
	111681.581866,
	111685.341207,
	111688.423454,
	111690.824758,
	111692.542121,
	111693.573398,
	111693.917300,
}

// Pre-computed meters per degree of longitude, 0 to 90 degrees.
var metersPerDegreeLng = []float64{
	111319.458000,
	111302.616974,
	111252.098718,
	111167.917694,
	111050.098004,
	110898.673387,
	110713.687212,
	110495.192475,
	110243.251793,
	109957.937389,
	109639.331091,
	109287.524313,
	108902.618048,
	108484.722850,
	108033.958821,
	107550.455593,
	107034.352306,
	106485.797595,
	105904.949560,
	105291.975746,
	104647.053118,
	103970.368032,
	103262.116205,
	102522.502685,
	101751.741816,
	100950.057205,
	100117.681681,
	99254.857257,
	98361.835087,
	97438.875421,
	96486.247557,
	95504.229794,
	94493.109377,
	93453.182442,
	92384.753961,
	91288.137676,
	90163.656041,
	89011.640152,
	87832.429679,
	86626.372793,
	85393.826090,
	84135.154513,
	82850.731268,
	81540.937739,
	80206.163399,
	78846.805721,
	77463.270075,
	76055.969635,
	74625.325274,
	73171.765457,
	71695.726129,
	70197.650606,
	68677.989454,
	67137.200370,
	65575.748052,
	63994.104080,
	62392.746775,
	60772.161066,
	59132.838355,
	57475.276364,
	55799.979000,
	54107.456194,
	52398.223755,
	50672.803207,
	48931.721632,
	47175.511501,
	45404.710512,
	43619.861416,
	41821.511839,
	40010.214113,
	38186.525088,
	36351.005951,
	34504.222041,
	32646.742657,
	30779.140870,
	28901.993324,
	27015.880042,
	25121.384226,
	23219.092055,
	21309.592482,
	19393.477028,
	17471.339574,
	15543.776155,
	13611.384743,
	11674.765042,
	9734.518271,
	7791.246950,
	5845.554682,
	3898.045944,
	1949.325862,
	0.000000,
}

// FastMetersPerDegreeLat returns meters per a latitude degree using a pre-computed lookup table.
func FastMetersPerDegreeLat(lat float64) float64 {
	i0 := int(math.Abs(math.Floor(lat)))
	i1 := int(math.Abs(math.Ceil(lat)))
	x := lat - float64(i0)
	m0 := metersPerDegreeLat[i0]
	m1 := metersPerDegreeLat[i1]
	return m0 + (m1-m0)*x
}

// FastMetersPerDegreeLng returns meters per a longitude degree using a pre-computed lookup table.
func FastMetersPerDegreeLng(lng float64) float64 {
	i0 := int(math.Abs(math.Floor(lng)))
	i1 := int(math.Abs(math.Ceil(lng)))
	x := lng - float64(i0)
	m0 := metersPerDegreeLng[i0]
	m1 := metersPerDegreeLng[i1]
	return m0 + (m1-m0)*x
}
