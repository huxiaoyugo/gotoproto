package generator



// 0 48
// a 97
// A 65
func lowerStr(str string) string {

	res := ""
	for _, item := range str {
		if item >= 'A' && item <= 'Z' {
			item += 32
		}
		res += string(item)
	}
	return res
}


// 驼峰命名改为下划线
func ToUnderLine(expr string) string {

	resByte := make([]byte, 0)
	for index, char := range []byte(expr) {
		if char >= 'A' && char <= 'Z' {
			char += 32
			if index != 0 {
				resByte = append(resByte, '_')
			}
		}
		resByte = append(resByte, char)
	}
	return string(resByte)
}
