/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';
import * as _ from 'lodash';
import { Table, Button, Row, Col, Card, Input, Icon, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';


class CountryManager extends React.Component {
    state = {
        townSelectedRowKeys: [],  // Check here to configure the default column
        countrySelectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        townsData: [],
        countriesData: [],
        selectedTown: '',
        selectedTownId: 0,
    };
    componentDidMount() {
        this.start();
    }
    start = () => {
        this.setState({ loading: true });
        this.fetchTownsData();
    };

    fetchTownsData = () => {
        const { fetchData } = this.props
        let tempTownId
        fetchData({funcName: 'fetchTowns', stateName: 'townsData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.towns === undefined) return
            tempTownId = res.data.towns[0].id
            this.setState({
                townsData: [...res.data.towns.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
                selectedTown: res.data.towns[0].name || '',
                selectedTownId: tempTownId || 0,
            });

            this.fetchCountriesData(tempTownId);
        });
    }

    fetchCountriesData = townId => {
        if (townId === undefined) return
        const { fetchData } = this.props
        fetchData({funcName: 'fetchCountries', stateName: 'countriesData', 
            params: {townId}}).then(res => {
            this.setState({
                countriesData: [...res.data.countries.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        }).catch(err => {
            this.setState({
                countriesData: [],
            })
        });
    }

    onTownSelectChange = (selectedRowKeys) => {
        if(selectedRowKeys.length > 0){
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length-1]]
        }
        
        this.setState({ townSelectedRowKeys: selectedRowKeys });
    };

    onTownRowClick = (record) => {
        const { townSelectedRowKeys } = this.state
        this.setState({
            selectedTown: record.name,
            selectedTownId: record.id,
            townSelectedRowKeys: townSelectedRowKeys.length > 0 && 
                townSelectedRowKeys[0] === record.id ? [] : [record.id],
        });
        this.fetchCountriesData(record.id)
    }

    onCountrySelectChange = (selectedRowKeys) => {
        if(selectedRowKeys.length > 0){
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length-1]]
        }
        
        this.setState({ countrySelectedRowKeys: selectedRowKeys });
    };

    onCountryRowClick = (record, index, event) => {
        const { countrySelectedRowKeys } = this.state
        this.setState({
            countrySelectedRowKeys: countrySelectedRowKeys.length > 0 &&
                countrySelectedRowKeys[0] === record.id ? [] : [record.id],
        })
    }

    handleAddTown = () => {
        let hasNewTown = false
        if(this.state.townsData && this.state.townsData[0] && this.state.townsData[0].key === -1) {
            hasNewTown = true
        }
        if(hasNewTown) return
        this.setState({
            townsData: [{
                key: -1,
                id: -1,
                name: '',
            }, ...this.state.townsData]
        });
    }

    handleCancelEditTown = () => {
        let tmpTownsData = [...this.state.townsData.filter(item => item.id !== -1)]
        this.setState({
            townsData: tmpTownsData,
            selectedTown: tmpTownsData[0].name,
            selectedTownId: tmpTownsData[0].id,
        })
    }

    handleDeleteTown = () => {
        const { fetchData } = this.props
        const { townSelectedRowKeys } = this.state
        if(townSelectedRowKeys.length === 0 || townSelectedRowKeys[0] === -1) return
        fetchData({funcName: 'deleteTown', params: {townId: townSelectedRowKeys[0]}, stateName: 'deleteTownStatus'})
            .then(res => {
                message.success('删除成功')
                this.fetchTownsData() 
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    onNewTownChange = (dataIndex, value) => {
        this.setState({
            [dataIndex]: value,
        })
    }


    onNewTownSave = () => {
        const { fetchData } = this.props
        const keys = _.keys(this.state)
        const PREFIX = 'town.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if(_.startsWith(key, PREFIX)){
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }
        fetchData({funcName: 'newTown', params: obj, stateName: 'newTownStatus'})
            .then(res => {
                message.success('创建成功')
                this.fetchTownsData() 
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    handleAddCountry = () => {
        let hasNewCountry = false
        if(this.state.countriesData && this.state.countriesData[0] && this.state.countriesData[0].key === -1) {
            hasNewCountry = true
        }
        if(hasNewCountry) return
        this.setState({
            countriesData: [{
                key: -1,
                id: -1,
                name: '',
            }, ...this.state.countriesData]
        });
    }

    handleCancelEditCountry = () => {
        let tmpCountriesData = [...this.state.countriesData.filter(item => item.id !== -1)]
        this.setState({
            countriesData: tmpCountriesData
        })
    }

    handleDeleteCountry = () => {
        const { fetchData } = this.props
        const { countrySelectedRowKeys, selectedTownId } = this.state
        if(countrySelectedRowKeys.length === 0 || countrySelectedRowKeys[0] === -1) return
        fetchData({funcName: 'deleteCountry', params: {countryId: countrySelectedRowKeys[0]}, stateName: 'deleteCountryStatus'})
            .then(res => {
                message.success('删除成功')
                this.fetchCountriesData(selectedTownId)
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    onNewCountryChange = (dataIndex, value) => {
        this.setState({
            [dataIndex]: value,
        })
    }

    onNewCountrySave = () => {
        const { fetchData } = this.props
        const { selectedTownId } = this.state
        const keys = _.keys(this.state)
        const PREFIX = 'country.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if(_.startsWith(key, PREFIX)){
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }
        obj.town_id = selectedTownId

        fetchData({funcName: 'newCountry', params: obj, stateName: 'newCountryStatus'})
            .then(res => {
                message.success('创建成功')
                this.fetchCountriesData(selectedTownId)
            }).catch(err => {
                let errRes = err.response
                if(errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    render() {

        const townColumns = [{
            title: '镇名',
            dataIndex: 'name',
            width: '50%',
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell dataIndex='town.name' value={record.name} onChange={this.onNewTownChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: '50%',
            render: (text, record) => {
                if (record.id === -1){
                    return <EditableCell type="opt" onSave={this.onNewTownSave} onCancel={this.handleCancelEditTown}/>
                }
                var createAt = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_LONG_FORMAT)
                return createAt;
            }
        }];
        
        const countryColumns = [{
            title: '村名',
            dataIndex: 'name',
            width: '50%',
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell dataIndex='country.name' value={record.name} onChange={this.onNewCountryChange}/>
                }else{
                    return <a>{text}</a>
                }
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: '50%',
            render: (text, record) => {
                if (record.id === -1){
                    return <EditableCell type="opt" onSave={this.onNewCountrySave} onCancel={this.handleCancelEditCountry}/>
                }
                var createAt = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_LONG_FORMAT)
                return createAt;
            }
        }];

        const { loading, townSelectedRowKeys, countrySelectedRowKeys, selectedTown,
            townsData, countriesData } = this.state;
        const townRowSelection = {
            selectedRowKeys: townSelectedRowKeys,
            onChange: this.onTownSelectChange,
            hideDefaultSelections: true,
            type: 'radio'
        };
        const countryRowSelection = {
            selectedRowKeys: countrySelectedRowKeys,
            onChange: this.onCountrySelectChange,
            ideDefaultSelections: true,
            type: 'radio'
        }

        const hasSelectedTown = townSelectedRowKeys.length > 0
        const hasSelectedCountry = countrySelectedRowKeys.length > 0
        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="村镇管理" second="" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={12}>
                        <div className="gutter-box">
                            <Card title="镇列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAddTown}
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteTown}
                                            disabled={!hasSelectedTown} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={townRowSelection} columns={townColumns} dataSource={townsData} 
                                    pagination={false} size="small"
                                    onRow={(record) => ({
                                        onClick: () => this.onTownRowClick(record),
                                    })}
                                />
                            </Card>
                        </div>
                    </Col>
                    <Col className="gutter-row" md={12}>
                        <div className="gutter-box">
                            <Card title={`${selectedTown + (selectedTown?"-":"")}村列表`} bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAddCountry}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteCountry}
                                            disabled={!hasSelectedCountry} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={countryRowSelection} columns={countryColumns} dataSource={countriesData} 
                                    pagination={false} size="small"
                                    onRow={(record) => ({
                                        onClick: () => this.onCountryRowClick(record),
                                    })}
                                />
                            </Card>
                        </div>
                    </Col>
                </Row>
            </div>
        );
    }
}

const mapStateToProps = state => {
    const { 
        townsData = {data: {count:0, towns:[]}}, 
        fetchCountries = {data: {count:0, countries:[]}} 
    } = state.httpData;
    return { townsData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(CountryManager)

