package views

import (
	"fmt"
	"testing"
)

func TestDiffCoverageInfo(t *testing.T) {
	configuration := ConfigurationJson{
		TestData: ConfigureJsonSubItems{ConfigureJsonSubItem{Version: "1.21.1.1"}, ConfigureJsonSubItem{Version: "1.123.1.1"}, ConfigureJsonSubItem{Version: "1.2.1.1"}},
		DiffData: ConfigureJsonSubItems{ConfigureJsonSubItem{Version: "11.21.1.1"}, ConfigureJsonSubItem{Version: "1.2.01.1"}, ConfigureJsonSubItem{Version: "1.2.1.1"}},
	}
	DiffCoverageInfo(&configuration)
}

func TestMergeLineList(t *testing.T) {
	lineList1 := "1 2 3 4 5 6 7 8 9"
	lineList2 := "3 5 6 7 8 9 10 11 12"

	expect := "1 2 3 4 5 6 7 8 9 10 11 12"
	got := MergeLineList(lineList1, lineList2)
	if expect != got {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
}

func TestMergeCoverageInfoForTwoIps(t *testing.T) {
	coverInfo1 := RedisDataFormat{
		CoverageInfo: UploadCoverage{
			TestId: "1",
			CoverageDataList: []UploadFileCoverageInfo{
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/c.go", RelPath: "/home/a/c.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "1 2 3 4",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "12 13 14 15",
						},
					},
				},
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/d.go", RelPath: "/home/a/d.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "1 2 3 4",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "12 13 14 15",
						},
					},
				},
			},
		},
	}
	coverInfo2 := RedisDataFormat{
		CoverageInfo: UploadCoverage{
			TestId: "2",
			CoverageDataList: []UploadFileCoverageInfo{
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/c.go", RelPath: "/home/a/c.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "5 6 7 8",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "16 17 18 19 20",
						},
					},
				},
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/d.go", RelPath: "/home/a/d.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "5 6 7 8",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "16 17 18 19 20",
						},
					},
				},
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/e.go", RelPath: "/home/a/e.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5",
							CoveredLineList: "1 2 3 4 5",
						},
					},
				},
			},
		},
	}

	expected := RedisDataFormat{
		CoverageInfo: UploadCoverage{
			TestId: "2",
			CoverageDataList: []UploadFileCoverageInfo{
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/c.go", RelPath: "/home/a/c.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "1 2 3 4 5 6 7 8",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "12 13 14 15 16 17 18 19 20",
						},
					},
				},
				UploadFileCoverageInfo{
					FileInfo: UploadFileInfo{
						ModuleName: "moduleName", Revision: "1234",
						SvnPath: "/home/a/d.go", RelPath: "/home/a/d.go",
					},
					FuncDetailList: []UploadFuncInfo{
						{
							FuncName: "Hello", EffectiveLineList: "1 2 3 4 5 6 7 8 9 10",
							CoveredLineList: "1 2 3 4 5 6 7 8",
						},
						{
							FuncName: "World", EffectiveLineList: "12 13 14 15 16 17 18 19 20",
							CoveredLineList: "12 13 14 15 16 17 18 19 20",
						},
					},
				},
			},
		},
	}
	got1 := MergeSameVersionCoverageInfoForTwoIps(coverInfo1, coverInfo2)
	got2 := MergeSameVersionCoverageInfoForTwoIps(coverInfo1, coverInfo2)

	fmt.Println("expect: ", expected)
	fmt.Println("got: ", got1)
	fmt.Println("got: ", got2)
}

func TestMergeCoverageInfoForAllIpsOfVersion(t *testing.T) {
	data := MergeCoverageInfoForAllIpsOfVersion("1.1.1.1")
	expect := "{2019-02-20 13:29:59.3928402 +0800 CST false {[{{  yes yes 1234 test1} [{say  1 2 3 4 5 6 7 8 9 10 1 3 4 5} {boodbye  12 13 14 15 16 17 18 19 20 13 14 15 16 17}]} {{  no no 1234 test2} [{Hello  1 2 3 4 5 6 7 8 9 10 1 3 4 5 6 7 8 9} {World  12 13 14 15 16 17 18 19 20 13 14 16 17 19 20}]}] 1 yudian sync}}"
	got := fmt.Sprintf("%v", *data)
	if got != expect {
		fmt.Printf("expect: %s\n, got: %s", expect, got)
	}
}

func TestMergeCoverageInfoForSameVersion(t *testing.T) {
	versionInfo := ConfigureJsonSubItem{
		Version: "1.1.1.1",
	}
	data := MergeCoverageInfoForSameVersion(&versionInfo)
	expect := "{2019-02-20 13:29:59.3928402 +0800 CST false {[{{  yes yes 1234 test1} [{say  1 2 3 4 5 6 7 8 9 10 1 3 4 5} {boodbye  12 13 14 15 16 17 18 19 20 13 14 15 16 17}]} {{  no no 1234 test2} [{Hello  1 2 3 4 5 6 7 8 9 10 1 3 4 5 6 7 8 9} {World  12 13 14 15 16 17 18 19 20 13 14 16 17 19 20}]}] 1 yudian sync}}"
	got := fmt.Sprintf("%v", *data)
	if got != expect {
		fmt.Printf("expect: %s\ngot: %s\n", expect, got)
	}

	versionInfo.IpList = []string{"192.168.1.1"}
	data = MergeCoverageInfoForSameVersion(&versionInfo)
	expect = "{2019-02-20 13:29:59.3928402 +0800 CST false {[{{  yes yes 1234 test1} [{say  1 2 3 4 5 6 7 8 9 10 1 3 4 5} {boodbye  12 13 14 15 16 17 18 19 20 13 14}]} {{  no no 1234 test2} [{Hello  1 2 3 4 5 6 7 8 9 10 1 3 4 5 6 7 8 9} {World  12 13 14 15 16 17 18 19 20 13 14 16 17 19}]}] 1 yudian sync}}"
	got = fmt.Sprintf("%v", *data)
	if got != expect {
		fmt.Printf("expect: %s\ngot: %s\n", expect, got)
	}
}

func Test_changeAnchorLineListToLineMap(t *testing.T) {
	anchorLineList := [][3]int{{98, 98, 0}, {99, 99, 0}, {100, 100, 0}, {101, 101, -1}, {102, 101, 1}, {102, 102, 0}, {103, 103, 0}, {104, 104, 0}, {112, 112, 0}, {113, 113, 0}, {114, 114, 0}, {115, 115, 1}, {115, 116, 1}, {115, 117, 1}, {115, 118, 0}, {116, 119, 0}, {117, 120, 0}, {133, 136, 0}, {134, 137, 0}, {135, 138, 0}, {136, 139, -1}, {137, 139, 0}, {138, 140, 0}, {139, 141, 0}, {140, 142, 0}}
	got := changeAnchorLineListToLineMap(anchorLineList, 150)
	fmt.Println("got: ", got)
	for i := 1; i < 150; i++ {
		fmt.Println("key: ", i, "value: ", got[i])
	}
}

func Test_changeLineFromOldToNewByDiffInfo(t *testing.T) {
	diffItem := SvnDiffItem{
		AnchorLineList: [][3]int{{98, 98, 0}, {99, 99, 0}, {100, 100, 0}, {101, 101, -1}, {102, 101, 1}, {102, 102, 0}, {103, 103, 0}, {104, 104, 0}, {112, 112, 0}, {113, 113, 0}, {114, 114, 0}, {115, 115, 1}, {115, 116, 1}, {115, 117, 1}, {115, 118, 0}, {116, 119, 0}, {117, 120, 0}, {133, 136, 0}, {134, 137, 0}, {135, 138, 0}, {136, 139, -1}, {137, 139, 0}, {138, 140, 0}, {139, 141, 0}, {140, 142, 0}},
	}
	lines := "114 115 116 117 118 119 120 135 136 137 138"
	expect := "114 118 119 120 121 122 123 138 139 140"
	got := changeLineFromOldToNewByDiffInfo(lines, &diffItem)
	if expect != got {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
}

func TestChangeCoverInfoFromOlderVersionToNewer(t *testing.T) {
	oldCoverageInfo := ReadCoverageFromStorage(1)
	got := ChangeCoverInfoFromOlderVersionToNewer(oldCoverageInfo, "1.1.2.2")
	fmt.Println(got)
}
