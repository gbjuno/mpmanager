/**
 * Created by Jingle on 2017/11/30.
 * 接口地址配置文件
 */

//easy-mock模拟数据接口地址
const EASY_MOCK = 'https://www.easy-mock.com/mock';
const MOCK_AUTH = EASY_MOCK + '/597b5ed9a1d30433d8411456/auth';         // 权限接口地址
export const MOCK_AUTH_ADMIN = MOCK_AUTH + '/admin';                           // 管理员权限接口
export const MOCK_AUTH_VISITOR = MOCK_AUTH + '/visitor';                       // 访问权限接口

const PROTOCOL = 'https'
const HOST = 'www.juntengshoes.cn'
const PORT = 80
const CONTEXT = 'api'
const VERSION = 'v1'

export const SERVER_ROOT = `${PROTOCOL}://${HOST}`
export const SERVER_URL = `${PROTOCOL}://${HOST}/${CONTEXT}/${VERSION}`


export const TOWN_URL = SERVER_URL + '/town'
export const TOWN_COUNTRY_URL = townId => `${TOWN_URL}/${townId}/country`
export const COUNTRY_URL = SERVER_URL + '/country'

export const COMPANY_URL = SERVER_URL + "/company"
export const USER_URL = SERVER_URL + "/user"
export const PLACE_URL = SERVER_URL + "/monitor_place"
export const PLACETYPE_URL = SERVER_URL + "/monitor_type"
export const SUMMARY_URL = SERVER_URL + "/summary"
export const PICTURE_URL = (filter) => SERVER_URL + `/monitor_place?day=${filter.day}&company_id=${filter.companyId}&pageNo=1&pageSize=1`