/**
 * Created by hao.cheng on 2017/4/16.
 */
import axios from 'axios';
import { get } from './tools';
import * as config from './config';

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


// 村镇管理API

export const fetchTowns = (filter={}) => {
    let url = config.TOWN_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const newTown = (town) => {
    return axios.post(config.TOWN_URL, {...town}, {headers: {Accept: 'application/json'}})
        .then(res => res.data).catch(err => console.log(err));
}

export const deleteTown = (town) => {
    if(town === undefined || town.townId === -1) return
    return axios.delete(config.TOWN_URL + "/" + town.townId)
        .then(res => res.data).catch(err => console.log(err));
}

export const fetchCountries = (filter={}) => {
    let url = config.TOWN_COUNTRY_URL(filter.townId)
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const newCountry = (country) => {
    return axios.post(config.COUNTRY_URL, {...country}, {headers: {Accept: 'application/json'}})
        .then(res => res.data).catch(err => console.log(err));
}

export const deleteCountry = (country) => {
    if(country === undefined || country.countryId === -1) return
    return axios.delete(config.COUNTRY_URL + "/" + country.countryId)
        .then(res => res.data).catch(err => console.log(err));
}



// 公司管理API

export const fetchCompanies = (filter={}) => {
    let url = `${config.COMPANY_URL}?pageNo=${filter.pageNo}&pageSize=${filter.pageSize}`
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

export const fetchCompaniesByCountryId = (filter={}) => {
    let url = `${config.COUNTRY_URL}/${filter.countryId}/company`
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}


// 用户管理API

export const fetchUsers = (filter={}) => {
    let url = config.USER_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}


// 地点管理API

export const fetchPlaces = (filter={}) => {
    let url = config.PLACE_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
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

// 图片管理API

export const fetchPictures = (filter={}) => {
    let url = config.PICTURE_URL
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}

// 效率低，查询很难用
// export const fetchPicturesByPlaceId = (filter={}) => {
//     let placeId = filter.placeId
//     let day = filter.day
//     let url = `${config.PLACE_URL}/${placeId}?scope=picture&day=${day}`
//     return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
// }

export const fetchPicturesWithPlace = (filter={}) => {
    let placeId = filter.placeId
    let day = filter.day
    let companyId = filter.companyId
    let url = config.PICTURE_URL({day, companyId})
    return axios.get(url ,{}).then(res => res.data).catch(err => console.log(err));
}