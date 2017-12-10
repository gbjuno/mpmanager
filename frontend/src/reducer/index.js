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
        case type.SEARCH_PICTURE:
            console.log('search picture reducer', action)
            return {
                ...state,
                [action.collection]: {...action.condition},
            };
        default:
            return {...state};
    }
}

export default combineReducers({
    httpData,
    searchFilter,
});
