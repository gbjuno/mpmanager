/**
 * Created by Jingle on 2017/12/11.
 */
import * as type from './type';
import * as http from '../axios/index';

const requestData = category => ({
    type: type.REQUEST_DATA,
    category
});
export const receiveData = (data, category) => ({
    type: type.RECEIVE_DATA,
    data,
    category
});
/**
 * 请求数据调用方法
 * @param funcName      请求接口的函数名
 * @param params        请求接口的参数
 */
export const fetchData = ({funcName, params, stateName}) => dispatch => {
    !stateName && (stateName = funcName);
    dispatch(requestData(stateName));
    return http[funcName](params).then(res => dispatch(receiveData(res, stateName)));
};

export const searchFilter = (scope, condition) => ({
    type: type.TRIGGER_SEARCH,
    scope,
    condition,
})

export const resetFilter = (scope) => ({
    type: type.RESET_SEARCH,
    scope,
})

export const updateMenu = (prevMenus, updateMenu, isNew, isSub) => ({
    type: type.UPDATE_WECHAT_MENU,
    prevMenus,
    updateMenu,
    isNew,
    isSub,
})

export const deleteMenu = (prevMenus, deleteMenu) => ({
    type: type.DELETE_WECHAT_MENU,
    prevMenus,
    deleteMenu,
})