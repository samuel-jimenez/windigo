/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

type Color uint32

func RGB(r, g, b byte) Color {
	return Color(uint32(r) | uint32(g)<<8 | uint32(b)<<16)
}

func (c Color) R() byte {
	return byte(c & 0xff)
}

func (c Color) G() byte {
	return byte((c >> 8) & 0xff)
}

func (c Color) B() byte {
	return byte((c >> 16) & 0xff)
}

func init() {

	Black = 0x000000
	DimGrey = 0x696969
	Grey = 0x808080
	DarkGrey = 0xa9a9a9
	Silver = 0xc0c0c0
	LightGrey = 0xd3d3d3
	Gainsboro = 0xdcdcdc
	WhiteSmoke = 0xf5f5f5
	White = 0xffffff

	RosyBrown = 0x8f8fbc
	IndianRed = 0x5c5ccd
	Brown = 0x2a2aa5
	FireBrick = 0x2222b2
	LightCoral = 0x8080f0
	Maroon = 0x000080
	DarkRed = 0x00008b

	Red = 0x0000ff
	Snow = 0xfafaff
	Salmon = 0x7280fa
	MistyRose = 0xe1e4ff
	Tomato = 0x4763ff
	DarkSalmon = 0x7a96e9
	OrangeRed = 0x0045ff
	Coral = 0x507fff
	LightSalmon = 0x7aa0ff
	Sienna = 0x2d52a0
	Chocolate = 0x1e69d2
	SaddleBrown = 0x13458b
	SeaShell = 0xeef5ff
	SandyBrown = 0x60a4f4
	PeachPuff = 0xb9daff
	Peru = 0x3f85cd
	Linen = 0xe6f0fa
	DarkOrange = 0x008cff
	Bisque = 0xc4e4ff
	Tan = 0x8cb4d2
	BurlyWood = 0x87b8de
	AntiqueWhite = 0xd7ebfa
	NavajoWhite = 0xaddeff
	BlanchedAlmond = 0xcdebff
	PapayaWhip = 0xd5efff
	Moccasin = 0xb5e4ff
	Wheat = 0xb3def5
	OldLace = 0xe6f5fd
	Orange = 0x00a5ff
	FloralWhite = 0xf0faff
	Goldenrod = 0x20a5da
	DarkGoldenrod = 0x0b86b8
	CornSilk = 0xdcf8ff
	Gold = 0x00d7ff
	Khaki = 0x8ce6f0
	LemonChiffon = 0xcdfaff
	PaleGoldenrod = 0xaae8ee
	DarkKhaki = 0x6bb7bd
	Beige = 0xdcf5f5
	LightGoldenrodYellow = 0xd2fafa
	Olive = 0x008080

	Yellow = 0x00ffff
	LightYellow = 0xe0ffff
	Ivory = 0xf0ffff
	OliveDrab = 0x238e6b
	YellowGreen = 0x32cd9a
	DarkOliveGreen = 0x2f6b55
	GreenYellow = 0x2fffad
	LawnGreen = 0x00fc7c
	Chartreuse = 0x00ff7f
	DarkSeaGreen = 0x8fbc8f
	ForestGreen = 0x228b22
	LimeGreen = 0x32cd32
	LightGreen = 0x90ee90
	PaleGreen = 0x98fb98
	DarkGreen = 0x006400
	Green = 0x008000

	Lime = 0x00ff00
	HoneyDew = 0xf0fff0
	SeaGreen = 0x578b2e
	MediumSeaGreen = 0x71b33c
	SpringGreen = 0x7fff00
	MintCream = 0xfafff5
	MediumSpringGreen = 0x9afa00
	MediumAquamarine = 0xaacd66
	Aquamarine = 0xd4ff7f
	Turquoise = 0xd0e040
	LightSeaGreen = 0xaab220
	MediumTurquoise = 0xccd148
	DarkSlateGrey = 0x4f4f2f
	PaleTurquoise = 0xeeeeaf
	Teal = 0x808000
	DarkCyan = 0x8b8b00

	Cyan = 0xffff00
	LightCyan = 0xffffe0
	Azure = 0xfffff0
	DarkTurquoise = 0xd1ce00
	CadetBlue = 0xa09e5f
	PowderBlue = 0xe6e0b0
	LightBlue = 0xe6d8ad
	DeepSkyBlue = 0xffbf00
	SkyBlue = 0xebce87
	LightSkyBlue = 0xface87
	SteelBlue = 0xb48246
	AliceBlue = 0xfff8f0
	SlateGrey = 0x908070
	LightSlateGrey = 0x998877
	DodgerBlue = 0xff901e
	LightSteelBlue = 0xdec4b0
	CornFlowerBlue = 0xed9564
	RoyalBlue = 0xe16941
	MidnightBlue = 0x701919
	Lavender = 0xfae6e6
	Navy = 0x800000
	DarkBlue = 0x8b0000
	MediumBlue = 0xcd0000

	Blue = 0xff0000
	GhostWhite = 0xfff8f8
	DarkSlateBlue = 0x8b3d48
	SlateBlue = 0xcd5a6a
	MediumSlateBlue = 0xee687b
	MediumPurple = 0xdb7093
	RebeccaPurple = 0x993366
	BlueViolet = 0xe22b8a
	Indigo = 0x82004b
	DarkOrchid = 0xcc3299
	DarkViolet = 0xd30094
	MediumOrchid = 0xd355ba
	Thistle = 0xd8bfd8
	Plum = 0xdda0dd
	Violet = 0xee82ee
	Purple = 0x800080
	DarkMagenta = 0x8b008b

	Magenta = 0xff00ff
	Orchid = 0xd670da
	MediumVioletRed = 0x8515c7
	DeepPink = 0x9314ff
	HotPink = 0xb469ff
	PaleVioletRed = 0x9370db
	LavenderBlush = 0xf5f0ff
	Crimson = 0x3c14dc
	Pink = 0xcbc0ff
	LightPink = 0xc1b6ff
}
