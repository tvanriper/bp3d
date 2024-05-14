package bp3d

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// Bin represents a container in which items will be put into.
type Bin struct {
	Name      string
	Width     float64
	Height    float64
	Depth     float64
	MaxWeight float64

	Items []*Item // Items that packed in this bin
}

type BinSlice []*Bin

func (bs BinSlice) Len() int { return len(bs) }
func (bs BinSlice) Less(i, j int) bool {
	return bs[i].GetVolume() < bs[j].GetVolume()
}
func (bs BinSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

type RevBinSlice []*Bin

func (bs RevBinSlice) Len() int { return len(bs) }
func (bs RevBinSlice) Less(i, j int) bool {
	return bs[j].GetVolume() < bs[i].GetVolume()
}
func (bs RevBinSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

// NewBin constructs new Bin with width w, height h, depth d, and max weight mw.
func NewBin(name string, w, h, d, mw float64) *Bin {
	return &Bin{
		Name:      name,
		Width:     w,
		Height:    h,
		Depth:     d,
		MaxWeight: mw,
		Items:     make([]*Item, 0),
	}
}

// GetName returns bin's name.
func (b *Bin) GetName() string {
	return b.Name
}

// GetWidth returns bin's width.
func (b *Bin) GetWidth() float64 {
	return b.Width
}

// GetHeight returns bin's height.
func (b *Bin) GetHeight() float64 {
	return b.Height
}

// GetDepth returns bin's depth.
func (b *Bin) GetDepth() float64 {
	return b.Depth
}

// GetDepth returns bin's volume.
func (b *Bin) GetVolume() float64 {
	return b.Width * b.Height * b.Depth
}

// GetUsedVolume returns the volume consumed by items in the bin.
func (b *Bin) GetUsedVolume() (result float64) {
	result = 0.0
	for _, item := range b.Items {
		result += item.GetVolume()
	}
	return
}

// GetAvailableVolume returns bin's available volume (after items added).
func (b *Bin) GetAvailableVolume() (result float64) {
	return b.GetVolume() - b.GetUsedVolume()
}

// GetVolumeUtilization returns the percentage of the space consumed.
func (b *Bin) GetVolumeUtilization() (result float64) {
	return (b.GetUsedVolume() * 100.0) / b.GetVolume()
}

// GetDepth returns bin's max weight.
func (b *Bin) GetMaxWeight() float64 {
	return b.MaxWeight
}

// PutItem tries to put item into pivot p of bin b.
func (b *Bin) PutItem(item *Item, p Pivot) (fit bool) {
	item.Position = p
	for i := 0; i < 6; i++ {
		item.RotationType = RotationType(i)
		d := item.GetDimension()
		if b.GetWidth() < p[0]+d[0] || b.GetHeight() < p[1]+d[1] || b.GetDepth() < p[2]+d[2] {
			continue
		}
		fit = true

		for _, ib := range b.Items {
			if ib.Intersect(item) {
				fit = false
				break
			}
		}

		if fit {
			b.Items = append(b.Items, item)
		}

		return
	}

	return
}

func (b *Bin) String() string {
	return fmt.Sprintf("%s(%vx%vx%v, max_weight:%v)", b.GetName(), b.GetWidth(), b.GetHeight(), b.GetDepth(), b.GetMaxWeight())
}

type RotationType int

const (
	RotationType_WHD RotationType = iota
	RotationType_HWD
	RotationType_HDW
	RotationType_DHW
	RotationType_DWH
	RotationType_WDH
)

var RotationTypeStrings = [...]string{
	"RotationType_WHD (w,h,d)",
	"RotationType_HWD (h,w,d)",
	"RotationType_HDW (h,d,w)",
	"RotationType_DHW (d,h,w)",
	"RotationType_DWH (d,w,h)",
	"RotationType_WDH (w,d,h)",
}

func (rt RotationType) String() string {
	return RotationTypeStrings[rt]
}

type Axis int

const (
	WidthAxis Axis = iota
	HeightAxis
	DepthAxis
)

type Pivot [3]float64
type Dimension [3]float64

func (pv Pivot) String() string {
	return fmt.Sprintf("%v,%v,%v", pv[0], pv[1], pv[2])
}

var startPosition = Pivot{0, 0, 0}

type Item struct {
	Name   string
	Width  float64
	Height float64
	Depth  float64
	Weight float64

	// Used during packer.Pack()
	RotationType RotationType
	Position     Pivot
}

type ItemSlice []*Item

func (is ItemSlice) Len() int { return len(is) }
func (is ItemSlice) Less(i, j int) bool {
	return is[i].GetVolume() > is[j].GetVolume()
}
func (is ItemSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

// NewItem returns an Item named name, with width w, height h, depth h, and
// weight w. The quantity defaults to one.
func NewItem(name string, w, h, d, wg float64) *Item {
	return &Item{
		Name:   name,
		Width:  w,
		Height: h,
		Depth:  d,
		Weight: wg,
	}
}

func (i *Item) GetName() string {
	return i.Name
}

func (i *Item) GetWidth() float64 {
	return i.Width
}

func (i *Item) GetHeight() float64 {
	return i.Height
}

func (i *Item) GetDepth() float64 {
	return i.Depth
}

func (i *Item) GetVolume() float64 {
	return i.Width * i.Height * i.Depth
}

func (i *Item) GetWeight() float64 {
	return i.Weight
}

func (i *Item) GetDimension() (d Dimension) {
	switch i.RotationType {
	case RotationType_WHD:
		d = Dimension{i.GetWidth(), i.GetHeight(), i.GetDepth()}
	case RotationType_HWD:
		d = Dimension{i.GetHeight(), i.GetWidth(), i.GetDepth()}
	case RotationType_HDW:
		d = Dimension{i.GetHeight(), i.GetDepth(), i.GetWidth()}
	case RotationType_DHW:
		d = Dimension{i.GetDepth(), i.GetHeight(), i.GetWidth()}
	case RotationType_DWH:
		d = Dimension{i.GetDepth(), i.GetWidth(), i.GetHeight()}
	case RotationType_WDH:
		d = Dimension{i.GetWidth(), i.GetDepth(), i.GetHeight()}
	}
	return
}

// Intersect checks whether there's an intersection between item i and item it.
func (i *Item) Intersect(i2 *Item) bool {
	return rectIntersect(i, i2, WidthAxis, HeightAxis) &&
		rectIntersect(i, i2, HeightAxis, DepthAxis) &&
		rectIntersect(i, i2, WidthAxis, DepthAxis)
}

// rectIntersect checks whether two rectangles from axis x and y of item i1 and i2
// has intersection or not.
func rectIntersect(i1, i2 *Item, x, y Axis) bool {
	d1 := i1.GetDimension()
	d2 := i2.GetDimension()

	cx1 := i1.Position[x] + d1[x]/2
	cy1 := i1.Position[y] + d1[y]/2
	cx2 := i2.Position[x] + d2[x]/2
	cy2 := i2.Position[y] + d2[y]/2

	ix := math.Max(cx1, cx2) - math.Min(cx1, cx2)
	iy := math.Max(cy1, cy2) - math.Min(cy1, cy2)

	// to avoid floating point errors, use loose criteria
	return 0.01 < (d1[x]+d2[x])/2-ix && 0.01 < (d1[y]+d2[y])/2-iy
}

func (i *Item) String() string {
	return fmt.Sprintf("%s(%vx%vx%v, weight: %v) pos(%s) rt(%s)", i.GetName(), i.GetWidth(), i.GetHeight(), i.GetDepth(), i.GetWeight(), i.Position, i.RotationType)
}

type Packer struct {
	FewestBoxes bool
	Bins        []*Bin
	Items       []*Item
	UnfitItems  []*Item // items that don't fit to any bin
}

func NewPacker() *Packer {
	return &Packer{
		FewestBoxes: false,
		Bins:        make([]*Bin, 0),
		Items:       make([]*Item, 0),
		UnfitItems:  make([]*Item, 0),
	}
}

func (p *Packer) AddBin(bins ...*Bin) {
	p.Bins = append(p.Bins, bins...)
}

func (p *Packer) AddItem(items ...*Item) {
	p.Items = append(p.Items, items...)
}

var (
	ErrInvalidBinsVolume = errors.New("invalid bins volume")
	ErrUnfitItemsExist   = errors.New("unfit items existing")
	ErrNoBins            = errors.New("no bins in packer")
	ErrNoItems           = errors.New("no items in packer")
)

func (p *Packer) Pack() error {
	sort.Sort(BinSlice(p.Bins))   // 昇順
	sort.Sort(ItemSlice(p.Items)) // 降順
	if len(p.Bins) == 0 {
		return ErrNoBins
	}
	if len(p.Items) == 0 {
		return ErrNoItems
	}

	maxVolumeItem := p.Items[0]
	maxVolumeBin := p.Bins[len(p.Bins)-1]
	if maxVolumeBin.GetVolume() < maxVolumeItem.GetVolume() {
		return ErrInvalidBinsVolume
	}

	itemVolumeSum := 0.0
	binVolumeSum := 0.0

	for _, item := range p.Items {
		itemVolumeSum += item.GetVolume()
	}
	for _, bin := range p.Bins {
		binVolumeSum += bin.GetVolume()
	}
	if binVolumeSum < itemVolumeSum {
		return ErrInvalidBinsVolume
	}

	if p.FewestBoxes {
		// Is there a bin that might hold all of the items?
		for _, bin := range p.Bins { // NOTE: sorted from smallest to largest.
			if bin.GetVolume() >= itemVolumeSum && len(p.Items) > 0 {
				// Yes... let's use it.
				p.Items = p.packToBin(bin, p.Items)
			}
		}
		for len(p.Items) > 0 {
			// No... so we want to attempt to pack with the
			// fewest possible boxes consuming the most volume.
			// How best to do that?
			// Calculate needed volume.
			need := 0.0
			for _, item := range p.Items {
				need += item.GetVolume()
			}
			// Find volume closest to this.
			found := 0.0
			var foundBin *Bin
			for _, bin := range p.Bins {
				if bin.GetAvailableVolume() >= found && bin.GetAvailableVolume() < need {
					foundBin = bin
					found = foundBin.GetAvailableVolume()
				}
			}
			if foundBin != nil {
				p.Items = p.packToBin(foundBin, p.Items)
			} else {
				// Nothing more we can do.
				break
			}
		}
	}

	for len(p.Items) > 0 {
		bin := p.FindFittedBin(p.Items[0])
		if bin == nil {
			p.unfitItem()
			continue
		}

		p.Items = p.packToBin(bin, p.Items)
	}

	if len(p.UnfitItems) > 0 {
		return ErrUnfitItemsExist
	}

	return nil
}

// unfitItem moves p.Items[0] to p.UnfitItems.
func (p *Packer) unfitItem() {
	if len(p.Items) == 0 {
		return
	}
	p.UnfitItems = append(p.UnfitItems, p.Items[0])
	p.Items = p.Items[1:]
}

// packToBin packs items to bin b. Returns unpacked items.
func (p *Packer) packToBin(b *Bin, items []*Item) (unpacked []*Item) {
	if !b.PutItem(items[0], startPosition) {

		if b2 := p.getBiggerBinThan(b); b2 != nil {
			return p.packToBin(b2, items)
		}

		return p.Items
	}

	// Pack unpacked items.
	for _, i := range items[1:] {
		var fitted bool
	lookup:

		// Try available pivots in current bin that are not intersect with
		// existing items in current bin.
		for pt := 0; pt < 3; pt++ {
			for _, ib := range b.Items {
				var pv Pivot
				switch Axis(pt) {
				case WidthAxis:
					pv = Pivot{ib.Position[0] + ib.GetWidth(), ib.Position[1], ib.Position[2]}
				case HeightAxis:
					pv = Pivot{ib.Position[0], ib.Position[1] + ib.GetHeight(), ib.Position[2]}
				case DepthAxis:
					pv = Pivot{ib.Position[0], ib.Position[1], ib.Position[2] + ib.GetDepth()}
				}

				if b.PutItem(i, pv) {
					fitted = true
					break lookup
				}
			}
		}

		if !fitted {
			for b2 := p.getBiggerBinThan(b); b2 != nil; b2 = p.getBiggerBinThan(b) {
				left := p.packToBin(b2, append(b2.Items, i))
				if len(left) == 0 {
					b = b2
					fitted = true
					break
				}
				b = b2
			}

			if !fitted {
				unpacked = append(unpacked, i)
			}
		}
	}

	return
}

func (p *Packer) getBiggerBinThan(b *Bin) *Bin {
	v := b.GetAvailableVolume()
	for _, b2 := range p.Bins {
		if b2.GetAvailableVolume() > v {
			return b2
		}
	}
	return nil
}

// FindFittedBin finds bin in which item i will be fitted into.
func (p *Packer) FindFittedBin(i *Item) *Bin {
	for _, b := range p.Bins {
		if !b.PutItem(i, startPosition) {
			continue
		}

		if len(b.Items) == 1 && b.Items[0] == i {
			// Clear items in bin as we previously just check whether item i
			// fits in bin b.
			b.Items = []*Item{}
		}

		return b
	}
	return nil
}
