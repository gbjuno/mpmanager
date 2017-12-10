/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Table, Button, Row, Col, Card, Input, Icon } from 'antd';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../BreadcrumbCustom';



class EditableCell extends React.Component {

    state = {
        value: this.props.value,
        editable: true,
    }
    handleChange = (e) => {
        const value = e.target.value;
        this.setState({ value });
    }
    check = () => {
        this.setState({ editable: false });
        if (this.props.onChange) {
            this.props.onChange(this.state.value);
        }
    }

    close = () => {
        this.setState({ editable: false });
        if (this.props.onCancel) {
            this.props.onCancel();
        }
    }

    edit = () => {
        this.setState({ editable: true });
    }

    render(){
        const { value, editable } = this.state;
        return (
            <div className="editable-cell">
                {
                editable ?
                    <div className="editable-cell-input-wrapper">
                    <Input
                        value={value}
                        onChange={this.handleChange}
                        onPressEnter={this.check}
                    />
                    <Icon
                        type="check"
                        className="editable-cell-icon-check"
                        onClick={this.check}
                    />
                    <Icon
                        type="close"
                        className="editable-cell-icon-close"
                        onClick={this.close}
                    />
                    </div>
                    :
                    <a target="_blank">{value}</a>
                }
            </div>
        )
    }
}

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
        console.log('selectedRowKeys changed: ', selectedRowKeys);
        if(selectedRowKeys.length > 0){
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length-1]]
        }
        
        this.setState({ townSelectedRowKeys: selectedRowKeys });
    };

    onTownRowClick = (record, index, event) => {
        console.log('select record...', record)
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
        console.log('selectedRowKeys changed: ', selectedRowKeys);
        if(selectedRowKeys.length > 0){
            selectedRowKeys = [selectedRowKeys[selectedRowKeys.length-1]]
        }
        
        this.setState({ countrySelectedRowKeys: selectedRowKeys });
    };

    onCountryRowClick = (record, index, event) => {
        console.log('select record...', record)
        const { countrySelectedRowKeys } = this.state
        this.setState({
            countrySelectedRowKeys: countrySelectedRowKeys.length > 0 &&
                countrySelectedRowKeys[0] === record.id ? [] : [record.id],
        })
    }

    handleAddTown = () => {
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
        console.log('tmpTownsData', tmpTownsData)
        this.setState({
            townsData: tmpTownsData,
            selectedTown: tmpTownsData[0].name,
            selectedTownId: tmpTownsData[0].id,
        })
    }

    handleDeleteTown = () => {
        const { fetchData } = this.props
        const { townSelectedRowKeys } = this.state
        if(townSelectedRowKeys.length === 0) return
        fetchData({funcName: 'deleteTown', params: {townId: townSelectedRowKeys[0]}, stateName: 'deleteTownStatus'})
            .then(res => {
                console.log('delete town successfully', res)
                this.fetchTownsData() 
            }).catch(err => console.log(err));
    }

    onNewTownSave = (name) => {
        const { fetchData } = this.props
        fetchData({funcName: 'newTown', params: {name}, stateName: 'newTownStatus'})
            .then(res => {
                console.log('create new town successfully', res)
                this.fetchTownsData() 
            }).catch(err => console.log(err));
        console.log('value--->', name)
    }

    handleAddCountry = () => {
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
        if(countrySelectedRowKeys.length === 0) return
        fetchData({funcName: 'deleteCountry', params: {countryId: countrySelectedRowKeys[0]}, stateName: 'deleteCountryStatus'})
            .then(res => {
                console.log('delete country successfully', res)
                this.fetchCountriesData(selectedTownId)
            }).catch(err => console.log(err))
    }

    onNewCountrySave = (name) => {
        const { fetchData } = this.props
        const { selectedTownId } = this.state
        fetchData({funcName: 'newCountry', params: {name, town_id: selectedTownId}, stateName: 'newCountryStatus'})
            .then(res => {
                console.log('create new coutry successfully', res)
                this.fetchCountriesData(selectedTownId)
            }).catch(err => {
                console.log(err)
                this.fetchCountriesData(selectedTownId)
            });
        console.log('value--->', name)
    }

    render() {

        const townColumns = [{
            title: '镇名',
            dataIndex: 'name',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.name} onChange={this.onNewTownSave} onCancel={this.handleCancelEditTown}/>
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: 80,
            render: (text, record) => {
                if (record.id === -1){
                    return ''
                }
                var createAt = new Date(text).toLocaleString('chinese',{hour12:false});
                return createAt;
            }
        }];
        
        const countryColumns = [{
            title: '村名',
            dataIndex: 'name',
            width: 40,
            render: (text, record) => {
                if(record.id === -1){
                    return <EditableCell value={record.name} onChange={this.onNewCountrySave} onCancel={this.handleCancelEditCountry}/>
                }else{
                    return <a href={record.url} target="_blank">{text}</a>
                }
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: 80,
            render: (text, record) => {
                if (record.id === -1){
                    return ''
                }
                var createAt = new Date(text).toLocaleString('chinese',{hour12:false});
                return createAt;
            }
        }];

        const { loading, townSelectedRowKeys, countrySelectedRowKeys, selectedTown,
            townsData, countriesData } = this.state;
        const townRowSelection = {
            selectedRowKeys: townSelectedRowKeys,
            onChange: this.onTownSelectChange,
        };
        const countryRowSelection = {
            selectedRowKeys: countrySelectedRowKeys,
            onChange: this.onCountrySelectChange,
        }
        return (
            <div className="gutter-example">
                <style>
                {`
                    .editable-cell {
                        position: relative;
                      }
                      
                      .editable-cell-input-wrapper,
                      .editable-cell-text-wrapper {
                        padding-right: 44px;
                      }
                      
                      .editable-cell-text-wrapper {
                        padding: 5px 24px 5px 5px;
                      }
                      
                      .editable-cell-icon,
                      .editable-cell-icon-check {
                        position: absolute;
                        right: 20px;
                        width: 20px;
                        cursor: pointer;
                      }

                      .editable-cell-icon-close {
                        position: absolute;
                        right: 0;
                        width: 20px;
                        cursor: pointer;
                      }
                      
                      .editable-cell-icon {
                        line-height: 18px;
                        display: none;
                      }
                      
                      .editable-cell-icon-check,
                      .editable-cell-icon-close {
                        line-height: 28px;
                      }
                      
                      .editable-cell:hover .editable-cell-icon {
                        display: inline-block;
                      }
                      
                      .editable-cell-icon:hover,
                      .editable-cell-icon-check:hover {
                        color: #108ee9;
                      }
                      
                      .editable-add-btn {
                        margin-bottom: 8px;
                      }
                `}
                </style>
                <BreadcrumbCustom first="安监管理" second="村镇管理" />
                <Row gutter={16}>
                    <Col className="gutter-row" md={12}>
                        <div className="gutter-box">
                            <Card title="镇列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAddTown}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteTown}
                                            disabled={loading} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={townRowSelection} columns={townColumns} dataSource={townsData} pagination={false}
                                        onRowClick={this.onTownRowClick}
                                />
                            </Card>
                        </div>
                    </Col>
                    <Col className="gutter-row" md={12}>
                        <div className="gutter-box">
                            <Card title={`${selectedTown + "-"}村列表`} bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAddCountry}
                                            disabled={loading} 
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleDeleteCountry}
                                            disabled={loading} 
                                    >删除</Button>
                                </div>
                                <Table rowSelection={countryRowSelection} columns={countryColumns} dataSource={countriesData} pagination={false}
                                        onRowClick={this.onCountryRowClick}/>
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

