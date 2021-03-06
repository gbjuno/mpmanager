/**
 * Created by Jingle on 2017/12/11.
 */
import { combineReducers } from 'redux';
import * as type from '../action/type';
import * as _ from 'lodash'

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
                [action.scope]: mergeCondition(state, action.scope, action.condition),
            };
        case type.RESET_SEARCH:
            return {
                ...state,
                [action.scope]: {},
            }
        default:
            return {...state};
    }
}

const wechatLocal = (state = {}, action) => {
    switch (action.type){
        case type.UPDATE_WECHAT_MENU:
            return {
                ...state,
                mergedMenus: mergeMenu(action.prevMenus, action.updateMenu, action.isNew, action.isSub),
            };
        case type.DELETE_WECHAT_MENU:
            return {
                ...state,
                mergedMenus: deleteMenu(action.prevMenus, action.deleteMenu)
            }
        case type.HANDLE_WECHAT_ARTICLE_ATTRIBUTE:
            return {
                ...state,
                article: mergeAttribute(state.article, action.attribute, action.value)
            }
        default:
            return {...state};
    }
}

const mergeCondition = (state, scope, condition) => _.merge(state[scope], condition)

const mergeAttribute = (prevArticle, attribute, value) => _.merge(prevArticle, {[attribute]: value})

const mergeMenu = (prevMenus, updateMenu, isNew, isSub) => {
    if(updateMenu === null){
        return prevMenus
    }
    if(prevMenus){

        if(isNew){
            if(isSub){
                prevMenus.button.filter(m => {
                    if(isOneOfSubMenus(updateMenu, m)){
                        return true
                    }else {
                        return false
                    }
                }).map(m => {
                    delete m['url']
                    if(!m.sub_button){
                        m.sub_button = []
                    }
                    m.sub_button.push({
                        "type": "view", 
                        "name": "子菜单名称", 
                        "url": "", 
                        "sub_button": [ ]
                    })
                })
                return prevMenus
            }else{
                prevMenus.button.push({
                    "type": "view", 
                    "name": "菜单名称", 
                    "url": "", 
                    "sub_button": [ ]
                })
            }

        }

        prevMenus.button.map(m => {
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
                return sm
            })
            return m
        })
        return prevMenus
    }
    return null
}

const deleteMenu = (prevMenus, deleteMenu) => {
    if(prevMenus){
        let dm_fk = deleteMenu.frontend_key
        if(dm_fk.toString().split('-').length > 1){
            prevMenus.button.map(m => {
                let index = _.findIndex(m.sub_button, {frontend_key: dm_fk})
                m.sub_button.splice(index, 1)
                return m
            })
            return prevMenus
        }else{
            let index = _.findIndex(prevMenus.button, {frontend_key: dm_fk})
            prevMenus.button.splice(index, 1)
            return prevMenus
        }
    }
}

const isOneOfSubMenus = (subMenu, menu) => {
    const sm_fk = subMenu.frontend_key
    const m_fk = menu.frontend_key

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
