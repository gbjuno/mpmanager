/**
 * Created by Jingle on 2017/12/11.
 */
import { combineReducers } from 'redux';
import * as type from '../action/type';

const handleData = (state = {isFetching: true, data: {}}, action) => {
    switch (action.type) {
        case type.REQUEST_DATA:
            return {...state, isFetching: true};
        case type.RECEIVE_DATA:
            return {...state, isFetching: false, data: action.data};
        default:
            return {...state};
    }
};
const httpData = (state = {}, action) => {
    switch (action.type) {
        case type.RECEIVE_DATA:
        case type.REQUEST_DATA:
            return {
                ...state,
                [action.category]: handleData(state[action.category], action)
            };
        default:
            return {...state};
    }
};

const searchFilter = (state = {}, action) => {
    switch (action.type){
        case type.TRIGGER_SEARCH:
            return {
                ...state,
                [action.scope]: {...action.condition},
            };
        default:
            return {...state};
    }
}

const wechatLocal = (state = {}, action) => {
    switch (action.type){
        case type.UPDATE_WECHAT_MENU:
            console.log('menu action...', action)
            return {
                ...state,
                mergedMenus: mergeMenu(action.prevMenus, action.updateMenu),
            };
        default:
            return {...state};
    }
}

const mergeMenu = (prevMenus, updateMenu) => {
    if(updateMenu === null){
        return prevMenus
    }
    return null
}

export default combineReducers({
    httpData,
    searchFilter,
    wechatLocal,
});
