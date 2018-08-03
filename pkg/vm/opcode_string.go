// Code generated by "stringer -trimprefix Op -type Opcode opcodes.go"; DO NOT EDIT.

package vm

import "strconv"

const _Opcode_name = "NopDropDrop2DupDup2SwapOverPickRollTuckRetFailZeroPush1Push2Push3Push4Push5Push6Push7Push8PushBOneNeg1PushTNowPushARandPushLAddSubMulDivModDivModMulDivNotNegIncDecEqGtLtIndexLenAppendExtendSliceFieldFieldLDefCallDecoEndDefIfZIfNZElseEndIfSumAvgMaxMinChoiceWChoiceSortLookup"

var _Opcode_map = map[Opcode]string{
	0:   _Opcode_name[0:3],
	1:   _Opcode_name[3:7],
	2:   _Opcode_name[7:12],
	5:   _Opcode_name[12:15],
	6:   _Opcode_name[15:19],
	9:   _Opcode_name[19:23],
	12:  _Opcode_name[23:27],
	13:  _Opcode_name[27:31],
	14:  _Opcode_name[31:35],
	15:  _Opcode_name[35:39],
	16:  _Opcode_name[39:42],
	17:  _Opcode_name[42:46],
	32:  _Opcode_name[46:50],
	33:  _Opcode_name[50:55],
	34:  _Opcode_name[55:60],
	35:  _Opcode_name[60:65],
	36:  _Opcode_name[65:70],
	37:  _Opcode_name[70:75],
	38:  _Opcode_name[75:80],
	39:  _Opcode_name[80:85],
	40:  _Opcode_name[85:90],
	41:  _Opcode_name[90:95],
	42:  _Opcode_name[95:98],
	43:  _Opcode_name[98:102],
	44:  _Opcode_name[102:107],
	45:  _Opcode_name[107:110],
	46:  _Opcode_name[110:115],
	47:  _Opcode_name[115:119],
	48:  _Opcode_name[119:124],
	64:  _Opcode_name[124:127],
	65:  _Opcode_name[127:130],
	66:  _Opcode_name[130:133],
	67:  _Opcode_name[133:136],
	68:  _Opcode_name[136:139],
	69:  _Opcode_name[139:145],
	70:  _Opcode_name[145:151],
	72:  _Opcode_name[151:154],
	73:  _Opcode_name[154:157],
	74:  _Opcode_name[157:160],
	75:  _Opcode_name[160:163],
	77:  _Opcode_name[163:165],
	78:  _Opcode_name[165:167],
	79:  _Opcode_name[167:169],
	80:  _Opcode_name[169:174],
	81:  _Opcode_name[174:177],
	82:  _Opcode_name[177:183],
	83:  _Opcode_name[183:189],
	84:  _Opcode_name[189:194],
	96:  _Opcode_name[194:199],
	112: _Opcode_name[199:205],
	128: _Opcode_name[205:208],
	129: _Opcode_name[208:212],
	130: _Opcode_name[212:216],
	136: _Opcode_name[216:222],
	137: _Opcode_name[222:225],
	138: _Opcode_name[225:229],
	142: _Opcode_name[229:233],
	143: _Opcode_name[233:238],
	144: _Opcode_name[238:241],
	145: _Opcode_name[241:244],
	146: _Opcode_name[244:247],
	147: _Opcode_name[247:250],
	148: _Opcode_name[250:256],
	149: _Opcode_name[256:263],
	150: _Opcode_name[263:267],
	151: _Opcode_name[267:273],
}

func (i Opcode) String() string {
	if str, ok := _Opcode_map[i]; ok {
		return str
	}
	return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
}
