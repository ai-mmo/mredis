package mredis

// SplitFromScore 将uint64分数拆分为两个uint64值
// 参数:
//   - score: 要拆分的分数值（uint64类型）
//
// 返回:
//   - uint64: 高32位的值
//   - uint64: 低32位的值（取反后的结果）
//
// 注意:
//   - 此函数用于将一个uint64分数拆分为两个部分
//   - 第一个返回值是score的高32位
//   - 第二个返回值是score低32位取反后的结果
func SplitFromScore(score uint64) (uint64, uint64) {
	first := score >> 32
	second := ^uint64(uint32(score&0xFFFFFFFF)) & 0xFFFFFFFF
	return first, second
}
