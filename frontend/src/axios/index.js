/**
 * Created by Jingle Chen on 2018/3/11.
 */
import axios from 'axios';
import { get } from './tools';
import * as config from './config';

import createHistory from 'history/createBrowserHistory'

export const history = createHistory()

axios.defaults.withCredentials = true;

/**
 * 全局拦截器，当session过期或者没有登陆信息时，跳到登录页面
 */
axios.interceptors.response.use(response => {
    console.log('response ---> ', response)
    return response;
}, error => {
    console.log('error ---> interceptor --->', error)
    if(error.response && error.response && error.response.status === 401){
        window.location.href = config.PAGE_CONTEXT;
    }
    return Promise.reject(error);
});

export const getPros = () => axios.post('http://api.xitu.io/resources/github', {
    category: "trending",
    period: "day",
    lang: "javascript",
    offset: 0,
    limit: 30
}).then(function (response) {
    return response.data;
}).catch(function (error) {
    console.log(error);
});

export const npmDependencies = () => axios.get('./npm.json').then(res => res.data).catch(err => console.log(err));

export const weibo = () => axios.get('./weibo.json').then(res => res.data).catch(err => console.log(err));

const GIT_OAUTH = 'https://github.com/login/oauth';
export const gitOauthLogin = () => axios.get(`${GIT_OAUTH}/authorize?client_id=792cdcd244e98dcd2dee&redirect_uri=http://localhost:3006/&scope=user&state=reactAdmin`);
export const gitOauthToken = code => axios.post('https://cors-anywhere.herokuapp.com/' + GIT_OAUTH + '/access_token', {...{client_id: '792cdcd244e98dcd2dee',
    client_secret: '81c4ff9df390d482b7c8b214a55cf24bf1f53059', redirect_uri: 'http://localhost:3006/', state: 'reactAdmin'}, code: code}, {headers: {Accept: 'application/json'}})
    .then(res => res.data).catch(err => console.log(err));
export const gitOauthInfo = access_token => axios({
    method: 'get',
    url: 'https://api.github.com/user?access_token=' + access_token,
}).then(res => res.data).catch(err => console.log(err));

// easy-mock数据交互
// 管理员权限获取
export const admin = () => get({url: config.MOCK_AUTH_ADMIN});

// 访问权限获取
export const guest = () => get({url: config.MOCK_AUTH_VISITOR});

// real login
export const authLogin = (authBody) => axios.request({
    method: 'post',
    url: config.LOGIN_URL,
    maxRedirects: 0,
    validateStatus: function(status) {
        return status >= 200 && status < 303;
    },
    data: authBody,
    headers:  {
        Accept: 'application/json',
    },
})


// 村镇管理API

export const fetchTowns = (filter={}) => {
    let url = config.TOWN_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const newTown = (town) => {
    return axios.post(config.TOWN_URL, {...town}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deleteTown = (town) => {
    if(town === undefined || town.townId === -1) return
    return axios.delete(config.TOWN_URL + "/" + town.townId)
        .then(res => res.data);
}

export const fetchCountries = (filter={}) => {
    let url = config.TOWN_COUNTRY_URL(filter.townId)
    return axios.get(url ,{}).then(res => res.data);
}

export const fetchCountriesWithoutTownId = (filter={}) => {
    return axios.get(config.COUNTRY_URL ,{}).then(res => res.data);
}

export const newCountry = (country) => {
    return axios.post(config.COUNTRY_URL, {...country}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deleteCountry = (country) => {
    if(country === undefined || country.countryId === -1) return
    return axios.delete(config.COUNTRY_URL + "/" + country.countryId)
        .then(res => res.data);
}



// 公司管理API

export const fetchCompanies = (filter={}) => {
    let url = `${config.COMPANY_URL}?pageNo=${filter.pageNo}&pageSize=${filter.pageSize}`
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const fetchCompaniesByCountryId = (filter={}) => {
    let url = `${config.COUNTRY_URL}/${filter.countryId}/company`
    return axios.get(url ,{}).then(res => res.data);
}

export const fetchUsersByCompanyId = (company) => {
    let url = `${config.COMPANY_URL}/${company.id}/user`
    return axios.get(url ,{}).then(res => res.data);
}

export const fetchPlacesByCompanyId = (company) => {
    let url = `${config.COMPANY_URL}/${company.id}/monitorplace`
    return axios.get(url ,{}).then(res => res.data);
}

export const newCompany = (company) => {
    return axios.post(config.COMPANY_URL, {...company}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const updateCompany = (company) => {
    return axios.put(config.COMPANY_URL  + "/" + company.id, {...company}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deleteCompany = (company) => {
    if(company === undefined || company.id === -1) return
    return axios.delete(config.COMPANY_URL + "/" + company.id)
        .then(res => res.data);
}

// 用户管理API

export const fetchUsers = (prevFilter={}) => {
    const filter = {
        pageNo: prevFilter.pageNo?prevFilter.pageNo:1,
        pageSize: prevFilter.pageSize?prevFilter.pageSize:10,
        name: prevFilter.name?prevFilter.name:'',
        phone: prevFilter.phone?prevFilter.phone:'',
    }
    let url = `${config.USER_URL}?pageNo=${filter.pageNo}&pageSize=${filter.pageSize}&name=${filter.name}&phone=${filter.phone}`
    return axios.get(url ,{}).then(res => res.data);
}


export const newUser = (user) => {
    return axios.post(config.USER_URL, {...user}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const updateUser = (user) => {
    return axios.put(config.USER_URL  + "/" + user.id, {...user}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deleteUser = (user) => {
    if(user === undefined || user.id === -1) return
    return axios.delete(config.USER_URL + "/" + user.id)
        .then(res => res.data);
}



// 地点管理API

export const fetchPlaces = (filter={}) => {
    let url = config.PLACE_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const newPlace = (place) => {
    return axios.post(config.PLACE_URL, {...place}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const updatePlace = (place) => {
    return axios.put(config.PLACE_URL  + "/" + place.id, {...place}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deletePlace = (place) => {
    if(place === undefined || place.id === -1) return
    return axios.delete(config.PLACE_URL + "/" + place.id)
        .then(res => res.data);
}

export const searchPlaces = (filter={}) => {
    let companyId = filter.companyId
    let url = config.SEARCH_PLACE_URL({companyId})
    return axios.get(url ,{}).then(res => res.data);
}

// 地点类型管理API

export const fetchPlaceTypes = (filter={}) => {
    let url = config.PLACETYPE_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}


// 统计报表API
export const fetchSummaries = (filter={}) => {
    let url = config.SUMMARY_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const searchSummaries = (filter={}) => {
    let url = config.SEARCH_SUMMARY_URL(filter)
    return axios.get(url ,{}).then(res => res.data);
}

// 图片管理API

export const fetchPictures = (filter={}) => {
    let url = config.PICTURE_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const updatePicture = (filter={}) => {
    let url = config.UPDATE_PICTURE_URL(filter)
    return axios.put(url , filter, {headers: {Accept: 'application/json'}}).then(res => res.data);
}

// 效率低，查询很难用
// export const fetchPicturesByPlaceId = (filter={}) => {
//     let placeId = filter.placeId
//     let day = filter.day
//     let url = `${config.PLACE_URL}/${placeId}?scope=picture&day=${day}`
//     return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
// }

export const fetchPicturesWithPlace = (filter={}) => {
    let day = filter.date
    let companyId = filter.companyId
    let url = config.PICTURE_URL({day, companyId})
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

// Wechat API

export const fetchMenus = () => {
    let url = config.WECHAT_MENU_URL
    return axios.get(url ,{}).then(res => res.data);
}


export const saveMenus = (payload) => {
    let url = config.WECHAT_MENU_URL
    return axios.post(url, payload, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}


export const fetchArticles = () => {
    let url = config.WECHAT_ARTICLE_URL
    return axios.get(url ,{}).then(res => res.data);
}


// 页面模板
export const fetchPages = (filter={}) => {
    let url = `${config.Page_URL}?pageNo=${filter.pageNo}&pageSize=${filter.pageSize}`
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const newPage= (page) => {
    return axios.post(config.Page_URL, {...page}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const updatePage= (page) => {
    return axios.put(config.Page_URL+ "/" + page.id, {...page}, {headers: {Accept: 'application/json'}})
        .then(res => res.data);
}

export const deletePage = (page) => {
    if(page === undefined || page.id === -1) return
    return axios.delete(config.Page_URL+ "/" + page.id)
        .then(res => res.data);
}