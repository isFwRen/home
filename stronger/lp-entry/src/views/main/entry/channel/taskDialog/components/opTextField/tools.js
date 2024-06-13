import moment from 'moment'
import { ignoreQestionMark, ignoreFreeValue, ignore } from '../../specificValidations/tools'

// 校验规则名
export const ruleKeyValues = {
  maxLen: 'max_length',
  minLen: 'min_length',
  fixLen: 'fixed_length',
  fixValue: 'included',
  maxVal: 'max_value',
  minVal: 'min_value',
  checkDate: 'check_date',
  specChar: 'free_value',

  '1': 'required',
  '2': 'spaced',
  '3': 'show_length',
  '4': 'digital',
  '5': 'amount',
  '6': 'chinese',
  '7': 'nonnegative',
  '8': 'alpha',
  '9': 'integer',
  '10': 'date',
  '11': 'id',
  '12': 'email',
  '13': 'phone',
  '14': 'name',
  
  includes: 'included',
  excludes: 'excluded',
  looseIncludes: 'loose_included',
  looseExcludes: 'loose_excluded'
}

export const MS_DAY = 86400000
export const MS_TODAY = new Date(moment().format('YYYY-MM-DD')).getTime()

// 身份证
export const ID_FIRST_TWO = ['11', '12', '13', '14', '15', '21', '22', '23', '31', '32', '33', '34', '35', '36', '37', '41', '42', '43', '44', '45', '46', '50', '51', '52', '53', '54', '61', '62', '63', '64', '65']
export const ID_LAST = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'x', 'X']

// 日期时限
export const CHECK_DATE_RULES = ['1', '2', '3', '4']

// 不需要校验的值(通用)
export const isGreenValue = ({ effectValidations, op, value }) => {
  // if(!value) {
  //   return true
  // }

  if(op !== 'opq') {

    if(ignore({ effectValidations, value })) return true

    // if(value.includes('?')) {
    //   return true
    // }
  }

  return false
}