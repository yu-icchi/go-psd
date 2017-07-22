package enginedata

type EngineData struct {
	EngineDict        EngineDict
	ResourceDict      ResourceDict
	DocumentResources DocumentResources
}

// EngineDict

type EngineDict struct {
	AntiAlias    int
	Editor       Editor
	GridInfo     GridInfo
	ParagraphRun ParagraphRun
	Rendered     Rendered
	StyleRun     StyleRun
}

type Editor struct {
	Text string
}

type GridInfo struct {
	AlignLineHeightToGridFlags bool
	GridColor                  GridColor
	GridIsOn                   bool
	GridLeading                int
	GridLeadingFillColor       GridLeadingFillColor
	GridSize                   int
	ShowGrid                   bool
}

type GridColor struct {
	Type   int
	Values []int
}

type GridLeadingFillColor struct {
	Type   int
	Values []int
}

type ParagraphRun struct {
	DefaultRunData DefaultRunData
	IsJoinable     int
	RunArray       []RunArray
	RunLengthArray []int
}

type DefaultRunData struct {
	Adjustments    Adjustments
	ParagraphSheet ParagraphSheet
}

type Adjustments struct {
	Axis []int
	XY   []int
}

type ParagraphSheet struct {
	DefaultStyleSheet int
	Properties        Properties
}

type RunArray struct {
	Adjustments    Adjustments
	ParagraphSheet ParagraphSheet
}

type Rendered struct {
	Shapes  Shapes
	Version int
}

type Shapes struct {
	Children         []Children
	WritingDirection int
}

type Children struct {
	Cookie     Cookie
	Lines      Lines
	Procession int
	ShapeType  int
}

type Cookie struct {
	Photoshop Photoshop
}

type Photoshop struct {
	Base      Base
	PointBase []int
	ShapeType int
}

type Base struct {
	ShapeType       int
	TransformPoint0 []int
	TransformPoint1 []int
	TransformPoint2 []int
}

type Lines struct {
	Children         []interface{}
	WritingDirection int
}

type StyleRun struct {
	DefaultRunData interface{} // todo
	IsJoinable     int
}

// ResourceDict

type ResourceDict struct {
}

// DocumentResources

type DocumentResources struct {
	FontSet                 []FontSet
	KinsokuSet              []KinsokuSet
	MojiKumiSet             []MojiKumiSet
	ParagraphSheetSet       []ParagraphSheetSet
	SmallCapSize            float64
	StyleSheetSet           []StyleSheetSet
	SubscriptPosition       float64
	SubscriptSize           float64
	SuperscriptPosition     float64
	SuperscriptSize         float64
	TheNormalParagraphSheet int
	TheNormalStyleSheet     int
}

type FontSet struct {
	FontType  int
	Name      string
	Script    int
	Synthetic int
}

type KinsokuSet struct {
	Hanging string
	Keep    string
	Name    string
	NoEnd   interface{} // todo
	NoStart string
}

type MojiKumiSet struct {
	InternalName string
}

type ParagraphSheetSet struct {
	DefaultStyleSheet int
	Name              string
	Properties        Properties
}

type Properties struct {
	AutoHyphenate      bool
	AutoLeading        float64
	Burasagari         bool
	ConsecutiveHyphens int
	EndIndent          int
	EveryLineComposer  bool
	FirstLineIndent    int
	GlyphSpacing       []int
	Hanging            bool
	HyphenatedWordSize int
	Justification      int
	KinsokuOrder       int
	LeadingType        int
	LetterSpacing      []int
	PostHyphen         int
	PreHyphen          int
	SpaceAfter         int
	SpaceBefore        int
	StartIndent        int
	WordSpacing        []float64
	Zone               int
}

type StyleSheetSet struct {
	Name           string
	StyleSheetData StyleSheetData
}

type StyleSheetData struct {
	AutoKerning        bool
	AutoLeading        bool
	BaselineDirection  int
	BaselineShift      int
	CharacterDirection int
	DLigatures         bool
	DiacriticPos       int
	FauxBold           bool
	FauxItalic         bool
	FillColor          FillColor
	FillFirst          bool
	FillFlag           bool
	Font               int
	FontBaseline       int
	FontCaps           int
	FontSize           int
	HindiNumbers       bool
	HorizontalScale    int
	Kashida            int
	Kerning            int
	Language           int
	Leading            int
	Ligatures          bool
	NoBreak            bool
	OutlineWidth       int
	Strikethrough      bool
	StrokeColor        StrokeColor
	StrokeFlag         bool
	StyleRunAlignment  int
	Tracking           int
	Tsume              int
	Underline          bool
	VerticalScale      int
	YUnderline         int
}

type FillColor struct {
	Type   int
	Values []int
}

type StrokeColor struct {
	Type   int
	Values []int
}
