/**
 * Created by Jingle on 2017/8/13.
 */
import React, { Component } from 'react';
// import { Router, Route, hashHistory, IndexRedirect } from 'react-router';
import { Route, Redirect, Switch } from 'react-router-dom';
import AuthBasic from '../components/auth/Basic';
import RouterEnter from '../components/auth/RouterEnter';
import PictureManager from '../containers/ux/PictureManager';
import CountryManager from '../containers/ux/CountryManager';
import CompanyManager from '../containers/ux/CompanyManager';
import PlaceManager from '../containers/ux/PlaceManager';
import UserManager from '../containers/ux/UserManager';
import SummaryManager from '../containers/ux/SummaryManager';
import PhotoStatus from '../containers/ux/PhotoStatus';
import MenuManager from '../containers/wechat/MenuManager';
import ArticleManager from '../containers/wechat/ArticleManager';
import ArticleForm from '../containers/wechat/ArticleForm';
import PageManager from '../containers/wechat/PageManager';


export default class CRouter extends Component {
    requireAuth = (permission, component) => {
        const { auth } = this.props;
        const { permissions } = auth.data;
        // const { auth } = store.getState().httpData;
        if (!permissions || !permissions.includes(permission)) return <Redirect to={'/login'} push />;
        return component;
    };
    render() {
        return (
            <Switch>

                {/* <Route exact path="/app/ux/tp" component={(props) => this.requireAuth('auth/ux/tp', <PictureManager {...props} />)} /> */}
                <Route exact path="/app/ux/tp" component={PictureManager} />
                <Route exact path="/app/ux/cz" component={CountryManager} />
                <Route exact path="/app/ux/gs" component={CompanyManager} />
                <Route exact path="/app/ux/dd" component={PlaceManager} />
                <Route exact path="/app/ux/yh" component={UserManager} />
                <Route exact path="/app/ux/tj" component={SummaryManager} />
                <Route exact path="/app/ux/wwc" component={PhotoStatus} />

                <Route exact path="/app/wechat/cd" component={MenuManager} />
                <Route exact path="/app/wechat/wz" component={ArticleManager} />
                <Route exact path="/app/wechat/wzbj" component={ArticleForm} />
                <Route exact path="/app/wechat/sc" component={PageManager} />
                <Route exact path="/app/wechat/xx" component={PageManager} />
                <Route exact path="/app/wechat/ym" component={PageManager} />

                <Route exact path="/app/auth/basic" component={AuthBasic} />
                <Route exact path="/app/auth/routerEnter" component={(props) => this.requireAuth('auth/testPage', <RouterEnter {...props} />)} />

                <Route render={() => <Redirect to="/404" />} />
            </Switch>
        )
    }
}