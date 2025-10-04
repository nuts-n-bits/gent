package main

import "fmt"

// -------- struct Revision --------

type Revision struct {
	Rid string // 1 = rid (optional)
	_fexists_1 bool
}

func (a *Revision) SetRid(b string) {
	a.Rid = b;
	a._fexists_1 = true;
}

func (a *Revision) GetRid() string {
	return a.Rid;
}

func (a *Revision) HasRid() bool {
	return a._fexists_1;
}

func (a *Revision) ToBytes() []byte {
	var fid uint64;
	var dat []byte;
	coll := make([]byte, 0);
	if a._fexists_1 {
		fid, dat = 1, []byte(a.Rid);
		coll = append(coll, encodeVarint(fid)...);
		coll = append(coll, encodeVarint(uint64(len(dat)))...);
		coll = append(coll, dat...);
	}
	return coll;
}

func (a *Revision) IsMissingRequired() (bool, []string) {
	coll := []string{};
	return len(coll) > 0, coll;
}

func (a *Revision) FromBytes(b []byte) {

}

// -------- struct Arsync4 --------

type Arsync4 struct {
	_fstore_2 string      // 2 = case_number (required)
	_fexist_2 bool        //
	_fstore_3 []Revision  // 3 = revs (repeated)
	_fstore_4 Revision    // 4 = latest_rev (required)
	_fexist_4 bool        //
	_fstore_5 []string    // 5 = addpc (repeated)
}

func (a *Arsync4) SetCaseNumber(b string) {
	a._fstore_2 = b;
	a._fexist_2 = true;
}

func (a *Arsync4) SetRevs(b []Revision) {
	a._fstore_3 = b;
}

func (a *Arsync4) SetLatestRev(b Revision) {
	a._fstore_4 = b;
	a._fexist_4 = true;
}

func (a *Arsync4) SetAddpc(b []string) {
	a._fstore_5 = b;
}

func (a *Arsync4) GetCaseNumber() string {
	return a._fstore_2;
}

func (a *Arsync4) GetRevs() []Revision {
	return a._fstore_3;
}

func (a *Arsync4) GetLatestRev() Revision {
	return a._fstore_4;
}

func (a *Arsync4) GetAddpc() []string {
	return a._fstore_5;
}

func (a *Arsync4) HasCaseNumber() bool {
	return a._fexist_2;
}

func (a *Arsync4) HasLatestRev() bool {
	return a._fexist_4;
}

func (a *Arsync4) IsMissingRequired() (bool, []string) {
	coll := []string{};
	if !a._fexist_2 {
		coll = append(coll, "case_number");
	}
	if !a._fexist_4 {
		coll = append(coll, "latest_rev");
	}
	for _, e := range a._fstore_3 {
		_, e := e.IsMissingRequired() 
		for _, e := range e {
			coll = append(coll, "revs/" + e);
		}
	}
	return len(coll) > 0, coll;
}

func (a *Arsync4) ToBytes() []byte {
	var fid uint64;
	var dat []byte;
	coll := make([]byte, 0);
	// 2 = case_number: string
	if a._fexist_2 {
		fid, dat = 2, []byte(a._fstore_2);
		coll = append(coll, encodeFrameHeader(fid, dat)...);
		coll = append(coll, dat...);
	}
	// 3 = revs: struct Revision (repeated)
	fid, dat = 3, make([]byte, 0);
	for _, e := range a._fstore_3 {
		t := e.ToBytes();
		dat = append(dat, encodeVarint(uint64(len(t)))...);
		dat = append(dat, t...);
	}
	if len(dat) != 0 {
		coll = append(coll, encodeFrameHeader(fid, dat)...);
		coll = append(coll, dat...);
	}
	// 4 = latest_revision: struct Revision
	if a._fexist_4 {
		fid, dat = 4, a._fstore_4.ToBytes();
		coll = append(coll, encodeFrameHeader(fid, dat)...);
		coll = append(coll, dat...);
	}
	// 5 = addpc: string (repeated)
	fid, dat = 5, make([]byte, 0);
	for _, e := range a._fstore_5 {
		t := []byte(e);
		dat = append(dat, encodeVarint(uint64(len(t)))...);
		dat = append(dat, t...);
	}
	if len(dat) != 0 {
		coll = append(coll, encodeFrameHeader(fid, dat)...);
		coll = append(coll, dat...);
	}
	return coll
}

func (a *Arsync4) FromBytes(b []byte) error {
	dwt := NewDefault()
	dwt.Load(b)
	for {
		if dwt.Remainder() == 0{
			break;
		}
		fid, data, signal, err := dwt.Step();
		if err != nil {
			return err;
		}
		if signal != 0 {
			return fmt.Errorf("incomplete message detected, signal = %d", signal);
		}
		switch fid {
		case 2: 
			a._fstore_2 = string(combine(data));
			a._fexist_2 = true;
		case 3: 
			d, err := DelimitCompactNoRem(data);
			if err != nil {
				 return err;
			}
			for _, e := range d {
				t := Revision{};
				t.FromBytes(combine(e));
				a._fstore_3 = append(a._fstore_3, t);
			}
		case 4:	
			t := Revision{};
			t.FromBytes(combine(data));
			a._fstore_4 = t;
			a._fexist_4 = true;
		}
	}
	return nil;
}

func combine(a [][]byte) []byte {
	coll := make([]byte, 0);
	for _, e := range a {
		coll = append(coll, e...);
	}
	return coll;
}

///// TSET //////

func p1(buffer []byte) {
	for i := range len(buffer) {
		fmt.Printf("0x%x ", buffer[i]);
	}
}

// func makeFrame(id uint64, data []byte) []byte {
// 	header := HeaderOf(id, len(data));
// 	combined := make([]byte, 0);
// 	combined = append(combined, header...);
// 	combined = append(combined, data...);
// 	return combined;
// }

// func main() {
// 	fmt.Println("starting...");
// 	arsync4 := Arsync4{};
// 	arsync4.SetCaseNumber("2025");
// 	arsync4.SetLatestRev(Revision{Rid: "33", _fexists_1: true});
// 	arsync4.SetAddpc([]string{"14", "15", "16"});
// 	buffer := arsync4.ToBytes();
// 	p1(buffer);
// 	arsync4_decode := Arsync4{};
// 	err := arsync4_decode.FromBytes(buffer);
// 	fmt.Printf("\nERR: %#v\n", err);
// 	fmt.Printf("CASENUM: %s\n", arsync4_decode.GetCaseNumber());
// 	fmt.Printf("ADDPC: %#v\n", arsync4_decode.GetAddpc());
// 	fmt.Printf("LATESTREV: %#v\n", arsync4_decode.GetLatestRev());
// 	fmt.Printf("REVS: %#v\n", arsync4_decode.GetRevs());
// 	{
// 		t1, t2 := arsync4_decode.IsMissingRequired();
// 		fmt.Printf("MISSING?: %#v, %#v", t1, t2);
// 	}
// }