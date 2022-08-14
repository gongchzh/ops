package al

type AppListS []AppList

func (a AppListS) Len() int {
	return len(a)
}
func (a AppListS) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a AppListS) Less(i, j int) bool {
	if a[j].AppName == a[i].AppName {
		return a[j].HostName > a[i].HostName
	}
	return a[j].AppName > a[i].AppName
}
