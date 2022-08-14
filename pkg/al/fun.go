package al

import (
	"ops/pkg/ctl"
	"ops/pkg/db"
)

type Funs []AppFun

func (funs Funs) GetFunIds() []int {
	var (
		ids []int
	)
	for _, v := range funs {
		ids = append(ids, v.FunId)
	}
	return ids
}

func (funs Funs) GetAllGroup() []AppGroup {
	var (
		gma []AppGroupSubMember
		gmm = make(map[int]map[int]int)
		gf  []AppGroupMember
		fm  = make(map[int]int)
	)
	if funs == nil || (len(funs) == 1 && funs[0].FunId == 0) {
		return nil
	}
	db.Db.SQL(`select b.* from app_group a,app_group_sub_member b
where a.app_group_id=b.app_group_id
and (a.state=100 or a.state=?)`, funs[0].AppZone).Find(&gma)
	db.Db.In("fun_id", funs.GetFunIds()).Find(&gf)
	for _, v := range gf {
		fm[v.AppGroupId] = 0
	}
	for _, v := range gma {
		if _, ok := gmm[v.SubGroupId]; !ok {
			gmm[v.SubGroupId] = make(map[int]int)
		}
		gmm[v.SubGroupId][v.AppGroupId] = 0
	}
	ctl.Debug(gmm)
	for _, v := range gf {
		GetParent(gmm, v.AppGroupId)
		ctl.Debug(gmm[v.AppGroupId])
	}
	return nil
}

type GroupParent struct {
	AppGroupId  int
	GroupParent []GroupParent
}

func GetParent(gmm map[int]map[int]int, grp int) map[int]map[int]int {
	for k, _ := range gmm[grp] {
		ctl.Debug(grp, k)
		if len(gmm[k]) == 0 {
			ctl.Debug(grp, k)
			continue
		}
		for s, _ := range gmm[k] {
			if _, ok := gmm[grp][s]; !ok {
				ctl.Debug(grp, s)
				gmm[grp][s] = 0
				GetParent(gmm, s)
			}
		}
	}
	return gmm
}
