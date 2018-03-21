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
                mergedMenus: mergeMenu(action.prevMenus, action.updateMenu, action.isNew,  action.isSub),
            };
        default:
            return {...state};
    }
}

const mergeMenu = (prevMenus, updateMenu, isNew, isSub) => {
    if(updateMenu === null){
        return prevMenus
    }
    if(prevMenus && prevMenus.menu){

        if(isNew){
            if(isSub){
                prevMenus.menu.button.filter(m => {
                    if(isOneOfSubMenus(updateMenu, m)){
                        return true
                    }else {
                        return false
                    }
                }).map(m => {
                    m.sub_button.push({
                        "type": "view", 
                        "name": "子菜单名称", 
                        "url": "", 
                        "sub_button": [ ]
                    })
                })
                return prevMenus
            }

        }

        prevMenus.menu.button.map(m => {
            if(m.frontend_key === updateMenu.frontend_key){
                m.name = updateMenu.name
                m.url = updateMenu.url
                return m
            }
            m.sub_button.map(sm => {
                if(sm.frontend_key === updateMenu.frontend_key){
                    sm.name = updateMenu.name
                    sm.url = updateMenu.url
                    return sm
                }
            })
        })
        return prevMenus
    }
    return null
}

const isOneOfSubMenus = (subMenu, menu) => {
    const sm_fk = subMenu.frontend_key
    const m_fk = menu.frontend_key

    console.log('xiangxinaiqing,,,,,', subMenu, menu)
    if(m_fk.toString() === sm_fk.split('-')[0]){
        return true
    }
    return false
}

export default combineReducers({
    httpData,
    searchFilter,
    wechatLocal,
});
