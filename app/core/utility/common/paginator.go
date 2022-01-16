package common

// Paging 分页器
//  paging:=Paging{}
//  paging.GetPages(page, pagesize, total)
//@author Zhiqing Guo
type Paging struct {
	Page  int64 `json:"page" form:"page"`   //当前页
	Size  int64 `json:"size" form:"size"`   //每页条数
	Total int64 `json:"total" form:"total"` //总条数
	Count int64 `json:"count" form:"count"` //总页数
	//StartNums int64 `json:"startnums" form:"startnums"` //起始条数
	//Nums      []int64  `json:"nums" form:"Nums"`//分页序数
	//NumsCount int64   `json:"numscount" form:"numscount"` //总页序数
}

// GetPages 获取分页信息
func (p *Paging) GetPages() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	var count int64
	if p.Total%p.Size == 0 {
		count = p.Total / p.Size
	} else {
		count = p.Total/p.Size + 1
	}
	if p.Page > count {
		p.Page = count
	}

	//p.StartNums = p.Size * (p.Page - 1)
	p.Count = count
}
