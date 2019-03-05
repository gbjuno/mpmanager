/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import { Table, Button, Row, Col, Calendar, Tabs, message, LocaleProvider, Card } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import VacationSearch from '../search/VacationSearch';
import * as config from '../../axios/config';
import * as utils from '../../utils';
import moment from 'moment';
import 'moment/locale/zh-cn';
import zhCN from 'antd/lib/locale-provider/zh_CN';
import vacon from '../../style/imgs/vacation_on.png';
import vacoff from '../../style/imgs/vacation_off.png';
import vaccom from '../../style/imgs/vacation_com.png';
import { fetchCompanyVacations } from '../../axios';
moment.locale('zh-cn');


const { Meta } = Card;
const TabPane = Tabs.TabPane;
const DIMESION = {
    FULL: 'finish_percentage_all',
    YEAR: 'finish_percentage_last_365_days',
    HALF_A_YEAR: 'finish_percentage_last_182_days',
    QUARTER: 'finish_percentage_last_90_days',
    MONTH: 'finish_percentage_last_30_days',
}


class VacationManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: 0,
        selectedRecord:{},
        currentPage: 1,
        visible: false,
        editable: false,
        hasNewRow: false,
        pageSize: 10,
        total: 0,
        expandedRowKeys: [],
        selectedVacationArray:[],
    };
    componentDidMount = () => {
        this.start();
    }

    start = () => {
        this.setState({ loading: true });
        this.fetchGlobalVacations();
    };

    fetchGlobalVacations = () => {
        const { fetchData } = this.props
        fetchData({
            funcName: 'fetchGlobalVacations', params: {
            }, stateName: 'globalVacations'
        }).then(res => {
            if (res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                total: res.data.count,
                loading: false,
            });
        });
    }

    fetchCompanyVacations = (companyId) => {
        const { fetchData } = this.props
        fetchData({
            funcName: 'fetchCompanyVacations', params: { companyId
            }, stateName: 'companyVacations'
        }).then(res => {
            if (res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                total: res.data.count,
                loading: false,
            });
        });
    }

    onNewRowChange = (dataIndex, value) => {
        this.setState({
            [dataIndex]: value,
        })
    }


    onSelectChange = (selectedRowKeys) => {
        if (selectedRowKeys.length > 0) {
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length - 1]]
        }

        this.setState({ selectedRowKeys });
    };

    onRowClick = (record, index, event) => {
        const { selectedRowKeys, editable } = this.state
        if (record.id === -1 || editable) {
            return
        }
        this.setState({
            selectedPeriodId: record.id,
            selectedRowKeys: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? [] : [record.id],
            selectedRecord: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? {} : record,
        }, () => {
            //this.fetchRelatedUserAndPlace(record.id);
        });

    }


    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        }, () => this.fetchData())
    }


    onPanelChange = (value, mode) => {
        console.log(value, mode);
    }

    onSelect = (moment) => {
        const { selectedVacationArray } = this.state
        const { globalVacations } = this.props
        if(this.isGlobalVactionRange(globalVacations, moment)){
            return
        }
        if(selectedVacationArray.length < 1){
            selectedVacationArray.push(moment)
        }else if(selectedVacationArray.length < 2){
            if(this.compareDate(selectedVacationArray[0], moment) === 0){
                selectedVacationArray.splice(0, 1)
            }else {
                selectedVacationArray.push(moment)
            }
        }else{
            if(this.compareDate(selectedVacationArray[0], moment) === 0){
                selectedVacationArray.splice(1, 1)
            }else{
                selectedVacationArray.splice(1, 1)
                selectedVacationArray.push(moment)
            }
        }
        this.setState({
            selectedVacationArray,
        })
    }

    dateCellRender = (moment) => {
        const { selectedVacationArray } = this.state
        const { globalVacations, companyVacations } = this.props
        if(this.isGlobalVactionRange(globalVacations, moment)){
            return (
                <div style={{textAlign: 'center'}}>
                    <img alt="假期" src={vacon} />
                    <div>假期</div>
                </div>
            )
        }
        if(this.isCompanyVactionRange(companyVacations, moment)){
            return (
                <div style={{textAlign: 'center'}}>
                    <img alt="公司假期" src={vaccom} />
                    <div>公司假期</div>
                </div>
            )
        }
        if(selectedVacationArray.length < 1){
           
            return
        } else if(selectedVacationArray.length < 2){
            if( this.compareDate(selectedVacationArray[0] ,moment) === 0){
                return (
                    <div style={{textAlign: 'center'}}>
                        <img alt="假期" src={vacoff} />
                        <div>设为假期</div>
                    </div>
                )
            }
        } else {
            if( this.compareDate(selectedVacationArray[0], selectedVacationArray[1]) > 0
                && this.compareDate(selectedVacationArray[0], moment) >= 0
                && this.compareDate(selectedVacationArray[1], moment) <= 0){
                return (
                    <div style={{textAlign: 'center'}}>
                        <img alt="假期" src={vacoff} />
                        <div>设为假期</div>
                    </div>
                )
            } else if( this.compareDate(selectedVacationArray[0], selectedVacationArray[1]) < 0
                && this.compareDate(selectedVacationArray[0], moment) <= 0
                && this.compareDate(selectedVacationArray[1], moment) >= 0){
                return (
                    <div style={{textAlign: 'center'}}>
                        <img alt="假期" src={vacoff} />
                        <div>设为假期</div>
                    </div>
                )
            }
        }
        
    }

    isInVaction = (moment, begin, end) => {
        return (this.compareDate(begin, moment) >= 0) && (this.compareDate(moment, end) >= 0)
    }

    compareDate = (a, b) => {
        return this.formatDate(b) - this.formatDate(a)
    }

    formatDate = (moment) => {
        return parseInt(moment.format('YYYYMMDD'))
    }

    formatDateEntry = (moment) => {
        return moment.format('YYYY-MM-DDThh:mm:ssZ')
    }

    convertToMoment = (text) => {
        return moment(new Date(text))
    }

    setVacation = () => {
        const { selectedVacationArray, selectedCompanyId } = this.state
        const { fetchData } = this.props
        
        let vacation = {}
        if(selectedVacationArray.length === 1){
            vacation ={
                start_at: this.formatDateEntry(selectedVacationArray[0]),
                end_at: this.formatDateEntry(selectedVacationArray[0]),
            }
        } else if(selectedVacationArray.length === 2){
            if(this.compareDate(selectedVacationArray[0], selectedVacationArray[1]) > 0){
                vacation = {
                    start_at: this.formatDateEntry(selectedVacationArray[0]),
                    end_at: this.formatDateEntry(selectedVacationArray[1]),
                }
            } else {
                vacation = {
                    start_at: this.formatDateEntry(selectedVacationArray[1]),
                    end_at: this.formatDateEntry(selectedVacationArray[0]),
                }
            }
        }
        if(selectedCompanyId && selectedCompanyId !== 0){
            vacation ={
                ...vacation,
                company_id: selectedCompanyId,
            }
            fetchData({funcName: 'createOrUpdateCompanyVacations', params: vacation, stateName: 'createOrUpdateCompanyVacationsStatus'})
            .then(res => {
                message.success('设置成功')
                this.fetchCompanyVacations() 
                this.setState({
                    selectedVacationArray: [],
                })
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
        } else {
            fetchData({funcName: 'createOrUpdateGlobalVacations', params: vacation, stateName: 'createOrUpdateGlobalVacationsStatus'})
            .then(res => {
                message.success('设置成功')
                this.fetchGlobalVacations() 
                this.setState({
                    selectedVacationArray: [],
                })
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
        }
    }

    cancelVacation = () => {
        const { selectedPeriodId, selectedCompanyId } = this.state
        const { fetchData } = this.props

        if(selectedCompanyId === 0) {
            fetchData({funcName: 'deleteGlobalVacations', params: {periodId: selectedPeriodId}, stateName: 'deleteGlobalVacationsStatus'})
                .then(res => {
                    message.success('删除成功')
                    this.fetchGlobalVacations() 
                    this.setState({
                        selectedVacationArray: [],
                    })
                }).catch(err => {
                    let errRes = err.response
                    if(errRes.data && errRes.data.status === 'error'){
                        message.error(errRes.data.error)
                    }
                });
        }else{
            fetchData({funcName: 'deleteCompanyVacations', params: {periodId: selectedPeriodId}, stateName: 'deleteCompanyVacationsStatus'})
                .then(res => {
                    message.success('删除成功')
                    this.fetchCompanyVacations(selectedCompanyId);
                    this.setState({
                        selectedVacationArray: [],
                    })
                }).catch(err => {
                    let errRes = err.response
                    if(errRes.data && errRes.data.status === 'error'){
                        message.error(errRes.data.error)
                    }
                });
        }
    }

    isGlobalVactionRange = (globalVacations, moment) => {
        let vacations = []
        if (!_.isEmpty(globalVacations) && !_.isEmpty(globalVacations.data)
            && !_.isEmpty(globalVacations.data.global_relax_periods)) {
            vacations = globalVacations.data.global_relax_periods;
        }
        let isVacation = false;
        vacations.forEach(vacation => {
            if (this.isInVaction(moment, this.convertToMoment(vacation.start_at), this.convertToMoment(vacation.end_at))) {
                isVacation = true;
                return;
            }
        });
        return isVacation;
    }

    isCompanyVactionRange = (companyVacations, moment) => {
        let vacations = []
        if (!_.isEmpty(companyVacations) && !_.isEmpty(companyVacations.data)
            && !_.isEmpty(companyVacations.data.company_relax_periods)) {
            vacations = companyVacations.data.company_relax_periods;
        }
        let isVacation = false;
        vacations.forEach(vacation => {
            if (this.isInVaction(moment, this.convertToMoment(vacation.start_at), this.convertToMoment(vacation.end_at))) {
                isVacation = true;
                return;
            }
        });
        return isVacation;
    }

    onCompanyChange = (value, companyName) => {

        if(value === 0){
            this.clearCompanyVacations()
        }else {
            this.fetchCompanyVacations(value)
        }
        this.setState({
            selectedCompanyId: value,
            selectedPeriodId: '',
            selectedRowKeys: [],
            selectedRecord: [],
        })
    }

    clearCompanyVacations = () => {

    }

    render() {
        const { loading, selectedRowKeys, selectedVacationArray, selectedCompanyId, selectedPeriodId } = this.state;
        const { companiesData, globalVacations, companyVacations } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        const vacationColumns = [{
            title: '假期从',
            dataIndex: 'start_at',
            width: '50%',
            render: (text, record) => {
                text = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                return <a>{text}</a>
            }
        },
        {
            title: '假期至',
            dataIndex: 'end_at',
            width: '50%',
            render: (text, record) => {
                text = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                return <a>{text}</a>
            }
        }];
        
        let vacationsWrappedData = []
        if(selectedCompanyId === 0){
            if (globalVacations.data && globalVacations.data.global_relax_periods) {
                vacationsWrappedData = [...globalVacations.data.global_relax_periods.map(item => { item.key = item.id; return item })]
            }
        }else{
            if (companyVacations.data && companyVacations.data.company_relax_periods) {
                vacationsWrappedData = [...companyVacations.data.company_relax_periods.map(item => { item.key = item.id; return item })]
            }
        }
        

        let options = [];

        const hasSelected = (selectedVacationArray.length > 0)
        let disabled = {}
        if(!hasSelected){
            disabled = {
                disabled: true,
            }
        }
        const hasSelectedVacation = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1


        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="假期管理" />
                <VacationSearch fetchData={fetchData} onChange={this.onCompanyChange} />
                <Row gutter={16}>
                <Col className="gutter-row" md={14}>
                <LocaleProvider locale={zhCN}>
                    <Calendar fullscreen onPanelChange={this.onPanelChange} onSelect={this.onSelect} dateCellRender={this.dateCellRender} />
                </LocaleProvider>
                </Col>
                <Col className="gutter-row" md={6}>
                    <Button type="primary" {...disabled} onClick={this.setVacation}>设置假期</Button>
                    <Button type="primary" disabled={!hasSelectedVacation} onClick={this.cancelVacation}>删除假期</Button>
                    <Table rowSelection={rowSelection} columns={vacationColumns} dataSource={vacationsWrappedData}
                        size="small"
                        style={{paddingTop: 10}}
                        onRow={(record) => ({
                            onClick: () => this.onRowClick(record),
                        })}
                        pagination={{
                            hideOnSinglePage: true,
                            onChange: this.handlePageChange,
                            defaultCurrent: 1,
                        }}
                    />
                </Col>
                </Row>
            </div>
        );
    }
}

const mapStateToProps = state => {
    const {
        companiesData = { data: { count: 0, towns: [] } },
        companyVacations = { data: { count: 0, users: [] } },
        globalVacations = { data: { count: 0, monitor_places: [] } },
    } = state.httpData;
    return { companiesData, companyVacations, globalVacations };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(VacationManager)

