/**
 * Created by Jingle Chen on 2017/12/7.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';
import * as _ from 'lodash'
import * as download from 'downloadjs'
import { Table, Button, Row, Col, Card, Input, Icon, Pagination, Modal, Upload, Tabs, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import { getPros } from '../../axios';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import EditableCell from '../../components/cells/EditableCell';
import CompanySearch from '../search/CompanySearch';
import * as config from '../../axios/config';

const TabPane = Tabs.TabPane;

class CompanyManager extends React.Component {
    state = {
        selectedRowKeys: [],  // Check here to configure the default column
        loading: false,
        companiesData: [],
        selectedCompany: '',
        selectedCompanyId: '',
        currentPage: 1,
        visible: false,
        editable: false,
        hasNewRow: false,
        pageSize: 10,
        total: 0,
        expandedRowKeys: [],
        // User table properties
        selectedUserKeys: [],
        selectedUserId: '',
        userEditable: false,
        hasNewUser: false,
        // Place table properties
        selectedPlaceKeys: [],
        selectedPlaceId: '',
        placeEditable: false,
        hasNewPlace: false,
    };
    componentDidMount = () => {
        this.start();
    }

    start = () => {
        this.setState({ loading: true });
        this.fetchData();
        this.fetchCountryListWithoutTownId();
        this.fetchPlaceType();
    };

    fetchData = () => {
        const { fetchData } = this.props
        const { currentPage, pageSize } = this.state
        let tempTownId
        fetchData({
            funcName: 'fetchCompanies', params: {
                pageNo: currentPage, pageSize: pageSize
            }, stateName: 'companiesData'
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

    fetchCountryListWithoutTownId = () => {
        const { fetchData } = this.props
        fetchData({ funcName: 'fetchCountriesWithoutTownId', stateName: 'countriesC2Data', params: {} }).then(res => {
            if (res === undefined || res.data === undefined || res.data.countries === undefined) return
            this.setState({
                countriesData: [...res.data.countries.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            });
        });
    }

    fetchPlaceType = () => {
        const { fetchData } = this.props
        fetchData({ funcName: 'fetchPlaceTypes', stateName: 'placeTypes' }).then(res => {
            if (res === undefined || res.data === undefined || res.data.monitor_types === undefined) return
            this.setState({
                placeTypes: [...res.data.monitor_types],
            }, () => {
            })
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
            selectedCompanyId: record.id,
            selectedCompany: record.name,
            selectedRowKeys: selectedRowKeys.length > 0 && selectedRowKeys[0] === record.id ? [] : [record.id],
        }, () => {
            this.fetchRelatedUserAndPlace(record.id);
        });

    }

    handleAdd = () => {
        this.setState({
            currentPage: 1,
            hasNewRow: true,
        });
    }

    handleCancelEditRow = () => {
        let tmpCompaniesData = [...this.state.companiesData.filter(item => item.id !== -1)]
        this.setState({
            editable: false,
            hasNewRow: false,
            companiesData: tmpCompaniesData,
            selectedRowKeys: [],
        })
    }

    handleModify = () => {
        this.setState({
            editable: true,
        })
    }

    handleDelete = () => {
        const { fetchData } = this.props
        const { selectedRowKeys, currentPage } = this.state
        if (selectedRowKeys.length === 0) return
        fetchData({
            funcName: 'deleteCompany', params: { id: selectedRowKeys[0] }, stateName: 'deleteCompanyStatus'
        }).then(res => {
            message.success('删除成功')
            this.setState({
                visible: false,
                selectedCompany: '',
                selectedCompanyId: '',
            })
            this.fetchData()
        }).catch(err => {
            let errRes = err.response
            if (errRes.data && errRes.data.status === 'error') {
                message.error(errRes.data.error)
                this.setState({
                    visible: false,
                })
            }
        });
    }

    /**
     * 准备删
     */
    showModal = () => {
        this.setState({
            visible: true,
        })
    }

    hideModal = () => {
        this.setState({
            visible: false,
        })
    }

    onRowSave = () => {
        const { editable } = this.state
        if (editable) {
            this.onUpdateRowSave()
        } else {
            this.onNewRowSave()
        }
    }

    onUpdateRowSave = () => {
        const { fetchData } = this.props
        const { selectedRowKeys } = this.state
        const keys = _.keys(this.state)
        const PREFIX = 'company.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }
        obj.id = selectedRowKeys[0]

        fetchData({ funcName: 'updateCompany', params: obj, stateName: 'updateCompanyStatus' })
            .then(res => {
                message.success('更新成功')
                this.fetchData()
                this.setState({
                    editable: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }

    onNewRowSave = () => {
        const { fetchData } = this.props
        const keys = _.keys(this.state)
        const PREFIX = 'company.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        fetchData({ funcName: 'newCompany', params: obj, stateName: 'newCompanyStatus' })
            .then(res => {
                message.success('创建成功')
                this.fetchData()
                this.setState({
                    editable: false,
                    hasNewRow: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }


    handlePageChange = (page, pageSize) => {
        this.setState({
            currentPage: page,
        }, () => this.fetchData())
    }

    uploadProps = () => {
        const props =
            {
                name: 'uploadFile',
                action: config.COMPANY_UPLOAD_URL, //TODO: 换成上传地址
                showUploadList: false,
                onChange(info) {
                    if (info.file.status !== 'uploading') {
                        console.log(info.file, info.fileList);
                    }
                    if (info.file.status === 'done') {
                        message.success(`${info.file.name}上传成功`);
                    } else if (info.file.status === 'error') {
                        message.error(`${info.file.name}上传失败`);
                    }
                },
            }
        return props
    }

    downloadFile = () => {
        let url = config.COMPANY_DOWNLOAD_URL; //TODO: 换成下载公司数据url,及相应的文件格式
        const x = new XMLHttpRequest;
        x.open("GET", url, true);
        x.responseType = "blob";
        x.onload = function (e) {
            download(x.response, "报表基础数据.xlsx", "application/octet-stream")
        }
        x.send();
    }



    /** TODO: BEGIN Additional Table here, maybe move to another component  */

    onUserChange = (selectedUserKeys) => {
        if (selectedUserKeys.length > 0) {
            selectedUserKeys = [selectedUserKeys[selectedUserKeys.length - 1]]
        }

        this.setState({ selectedUserKeys });
    };

    onUserSelect = (record, selected, selectedUserKeys) => {
        if (selected && selectedUserKeys[0] === record.id) {
            this.setState({
                selectedUserKeys: []
            })
        } else {

        }
    }

    onUserClick = (record, index, event) => {
        const { selectedUserKeys, userEditable } = this.state
        if (record.id === -1 || userEditable) {
            return
        }
        this.setState({
            selectedUserId: record.id,
            selectedUserKeys: selectedUserKeys.length > 0 && selectedUserKeys[0] === record.id ? [] : [record.id],
        });

    }

    handleAddUser = () => {
        this.setState({
            hasNewUser: true,
        });
    }

    handleCancelEditUser = () => {
        this.setState({
            userEditable: false,
            hasNewUser: false,
            selectedUserKeys: [],
        })
    }

    handleModifyUser = () => {
        this.setState({
            userEditable: true,
        })
    }

    onUserSave = () => {
        const { fetchData } = this.props
        const { userEditable, selectedRowKeys, selectedCompanyId, selectedUserKeys } = this.state
        let funcName, stateName, successMessage;

        const keys = _.keys(this.state)
        const PREFIX = 'user.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        obj.company_id = selectedCompanyId
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        if (userEditable) {
            funcName = 'updateUser'
            stateName = 'updateUserStatus'
            successMessage = '更新成功'
            obj.id = selectedUserKeys[0]
        } else {
            funcName = 'newUser'
            stateName = 'newUserStatus'
            successMessage = '创建成功'
        }

        fetchData({ funcName, params: obj, stateName })
            .then(res => {
                message.success(successMessage)
                this.fetchRelatedUserAndPlace(selectedCompanyId)
                this.setState({
                    userEditable: false,
                    hasNewUser: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }

    handleDeleteUser = () => {
        const { fetchData } = this.props
        const { selectedUserKeys, selectedCompanyId } = this.state
        if (selectedUserKeys.length === 0) return
        fetchData({
            funcName: 'deleteUser', params: { id: selectedUserKeys[0] }, stateName: 'deleteUserStatus'
        }).then(res => {
            message.success('删除成功')
            this.setState({
                visible: false,
                selectedUserKeys: [],
                selectedUserId: '',
            })
            this.fetchRelatedUserAndPlace(selectedCompanyId)
        }).catch(err => {
            let errRes = err.response
            if (errRes.data && errRes.data.status === 'error') {
                message.error(errRes.data.error)
                this.setState({
                    visible: false,
                })
            }
        });
    }

    // The upper is user operations

    onPlaceChange = (selectedPlaceKeys) => {
        if (selectedPlaceKeys.length > 0) {
            selectedPlaceKeys = [selectedPlaceKeys[selectedPlaceKeys.length - 1]]
        }

        this.setState({ selectedPlaceKeys });
    };

    onPlaceSelect = (record, selected, selectedPlaceKeys) => {
        if (selected && selectedPlaceKeys[0] === record.id) {
            this.setState({
                selectedPlaceKeys: []
            })
        } else {

        }
    }

    onPlaceClick = (record, index, event) => {
        const { selectedPlaceKeys, placeEditable } = this.state
        if (record.id === -1 || placeEditable) {
            return
        }
        this.setState({
            selectedPlaceId: record.id,
            selectedPlaceKeys: selectedPlaceKeys.length > 0 && selectedPlaceKeys[0] === record.id ? [] : [record.id],
        });

    }

    handleAddPlace = () => {
        this.setState({
            hasNewPlace: true,
        });
    }

    handleCancelEditPlace = () => {
        this.setState({
            placeEditable: false,
            hasNewPlace: false,
            selectedPlaceKeys: [],
        })
    }

    handleModifyPlace = () => {
        this.setState({
            placeEditable: true,
        })
    }

    onPlaceSave = () => {
        const { fetchData } = this.props
        const { placeEditable, selectedRowKeys, selectedCompanyId, selectedPlaceKeys } = this.state
        let funcName, stateName, successMessage;

        const keys = _.keys(this.state)
        const PREFIX = 'place.'
        const PREFIX_LEN = PREFIX.length;
        let obj = {}
        obj.company_id = selectedCompanyId
        for (let key of keys) {
            if (_.startsWith(key, PREFIX)) {
                let field = key.substring(PREFIX_LEN)
                obj[field] = this.state[key]
            }
        }

        if (placeEditable) {
            funcName = 'updatePlace'
            stateName = 'updatePlaceStatus'
            successMessage = '更新成功'
            obj.id = selectedPlaceKeys[0]
        } else {
            funcName = 'newPlace'
            stateName = 'newPlaceStatus'
            successMessage = '创建成功'
        }

        fetchData({ funcName, params: obj, stateName })
            .then(res => {
                message.success(successMessage)
                this.fetchRelatedUserAndPlace(selectedCompanyId)
                this.setState({
                    placeEditable: false,
                    hasNewPlace: false,
                })
            }).catch(err => {
                let errRes = err.response
                if (errRes.data && errRes.data.status === 'error') {
                    message.error(errRes.data.error)
                }
            });
    }

    handleDeletePlace = () => {
        const { fetchData } = this.props
        const { selectedPlaceKeys, selectedCompanyId } = this.state
        if (selectedPlaceKeys.length === 0) return
        fetchData({
            funcName: 'deletePlace', params: { id: selectedPlaceKeys[0] }, stateName: 'deletePlaceStatus'
        }).then(res => {
            message.success('删除成功')
            this.setState({
                visible: false,
                selectedPlaceKeys: [],
                selectedPlaceId: '',
            })
            this.fetchRelatedUserAndPlace(selectedCompanyId)
        }).catch(err => {
            let errRes = err.response
            if (errRes.data && errRes.data.status === 'error') {
                message.error(errRes.data.error)
                this.setState({
                    visible: false,
                })
            }
        });
    }

    // The upper is place operations


    additionalTable = () => {
        const { selectedUserKeys, userEditable, hasNewUser, selectedCompanyId,
            selectedPlaceKeys, placeEditable, hasNewPlace, placeTypes } = this.state
        const { usersInCompany, placesInCompany } = this.props
        const userRowSelection = {
            selectedRowKeys: selectedUserKeys,
            onChange: this.onUserChange,
            onSelect: this.onUserSelect,
            type: 'radio',
        }
        const placeRowSelection = {
            selectedRowKeys: selectedPlaceKeys,
            onChange: this.onPlaceChange,
            onSelect: this.onPlaceSelect,
            type: 'radio',
        }

        const userColumns = [
            {
                title: '用户名',
                dataIndex: 'name',
                key: 'name',
                width: '25%',
                render: (text, record) => {
                    if (record.id === -1 || (userEditable && record.id === selectedUserKeys[0])) {
                        return <EditableCell dataIndex='user.name' value={record.name} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
            {
                title: '手机',
                dataIndex: 'phone',
                key: 'phone',
                width: '25%',
                render: (text, record) => {
                    if (record.id === -1 || (userEditable && record.id === selectedUserKeys[0])) {
                        return <EditableCell dataIndex='user.phone' value={record.phone} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
            {
                title: '职位',
                dataIndex: 'job',
                key: 'job',
                width: '25%',
                render: (text, record) => {
                    if (record.id === -1 || (userEditable && record.id === selectedUserKeys[0])) {
                        return <EditableCell dataIndex='user.job' value={record.job} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
            {
                title: '创建时间',
                dataIndex: 'create_at',
                key: 'create_at',
                width: '25%',
                render: (text, record) => {
                    if (record.id === -1 || (userEditable && record.id === selectedUserKeys[0])) {
                        return <EditableCell type="opt" onSave={this.onUserSave} onCancel={this.handleCancelEditUser} />
                    }
                    var createAt = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                    return createAt;
                }
            },
        ];

        const placeColumns = [
            {
                title: '地点名',
                dataIndex: 'name',
                key: 'name',
                width: '33%',
                render: (text, record) => {
                    if (record.id === -1 || (placeEditable && record.id === selectedPlaceKeys[0])) {
                        return <EditableCell dataIndex='place.name' value={record.name} onChange={this.onNewRowChange} />
                    }
                    return text
                }
            },
            {
                title: '类型',
                dataIndex: 'monitor_type_name',
                key: 'monitor_type_name',
                width: '33%',
                render: (text, record) => {
                    if (record.id === -1 || (placeEditable && record.id === selectedPlaceKeys[0])) {
                        return <EditableCell dataIndex='place.monitor_type_id' value={record.monitor_type_id} onChange={this.onNewRowChange}
                            editType="select" valueType="int" options={placeTypes} placeholder="请选择类型" />
                    }
                    return text
                }
            },
            {
                title: '创建时间',
                dataIndex: 'create_at',
                key: 'create_at',
                width: '34%',
                render: (text, record) => {
                    if (record.id === -1 || (placeEditable && record.id === selectedPlaceKeys[0])) {
                        return <EditableCell type="opt" onSave={this.onPlaceSave} onCancel={this.handleCancelEditPlace} />
                    }
                    var createAt = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                    return createAt;
                }
            },
        ];

        let usersData = []
        if (usersInCompany.data && usersInCompany.data.users) {
            usersData = [...usersInCompany.data.users.map(item => { item.key = item.id; return item })]
        }
        if (hasNewUser) {
            usersData = [{
                key: -1,
                id: -1,
                name: '',
                phone: '',
                job: '',
                company_id: '',
            }, ...usersData]
        } else {
            usersData = [...usersData.filter(item => item.id !== -1)]
        }
        const hasSelectedUser = selectedUserKeys.length > 0 && selectedUserKeys[0] !== -1

        let placesData = []
        if (placesInCompany.data && placesInCompany.data.monitor_places) {
            placesData = [...placesInCompany.data.monitor_places.map(item => { item.key = item.id; return item })]
        }
        if (hasNewPlace) {
            placesData = [{
                key: -1,
                id: -1,
                name: '',
                monitor_type_id: '',
                company_id: '',
            }, ...placesData]
        } else {
            placesData = [...placesData.filter(item => item.id !== -1)]
        }
        const hasSelectedPlace = selectedPlaceKeys.length > 0 && selectedPlaceKeys[0] !== -1

        const hasSelectedCompany = selectedCompanyId !== undefined && selectedCompanyId !== ''

        return (
            <Tabs defaultActiveKey="1">
                <TabPane tab="用户" key="1">
                    <div style={{ marginBottom: 16 }}>
                        <Button type="primary" onClick={this.handleAddUser}
                            disabled={!hasSelectedCompany}
                        >新增</Button>
                        <Button type="primary" onClick={this.handleModifyUser}
                            disabled={!hasSelectedUser}
                        >修改</Button>
                        <Button type="primary" onClick={this.handleDeleteUser}
                            disabled={!hasSelectedUser}
                        >删除</Button>
                        <Modal
                            title="警告"
                            visible={this.state.userVisible}
                            onOk={this.hideModal}
                            onCancel={this.hideModal}
                            okText="确认"
                            cancelText="取消"
                        >
                            <p>确认删除</p>
                        </Modal>
                    </div>
                    <Table size="small" columns={userColumns} dataSource={usersData}
                        rowSelection={userRowSelection} pagination={false}
                        onRow={(record) => ({
                            onClick: () => this.onUserClick(record),
                        })}
                    />
                </TabPane>
                <TabPane tab="地点" key="2">
                    <div style={{ marginBottom: 16 }}>
                        <Button type="primary" onClick={this.handleAddPlace}
                            disabled={!hasSelectedCompany}
                        >新增</Button>
                        <Button type="primary" onClick={this.handleModifyPlace}
                            disabled={!hasSelectedPlace}
                        >修改</Button>
                        <Button type="primary" onClick={this.handleDeletePlace}
                            disabled={!hasSelectedPlace}
                        >删除</Button>
                        <Modal
                            title="警告"
                            visible={this.state.placeVisible}
                            onOk={this.hideModal}
                            onCancel={this.hideModal}
                            okText="确认"
                            cancelText="取消"
                        >
                            <p>确认删除</p>
                        </Modal>
                    </div>
                    <Table size="small" columns={placeColumns} dataSource={placesData}
                        rowSelection={placeRowSelection} pagination={false}
                        onRow={(record) => ({
                            onClick: () => this.onPlaceClick(record),
                        })}
                    />
                </TabPane>
            </Tabs>
        )
    }

    /** TODO: END Additional Table here, maybe move to another component  */

    fetchRelatedUserAndPlace = (selectedCompanyId) => {
        const { fetchData } = this.props
        fetchData({
            funcName: 'fetchUsersByCompanyId', params: {
                id: selectedCompanyId,
            }, stateName: 'usersInCompany'
        })
        fetchData({
            funcName: 'fetchPlacesByCompanyId', params: {
                id: selectedCompanyId,
            }, stateName: 'placesInCompany'
        })

    }

    render() {
        const { loading, selectedRowKeys, selectedTown, editable, hasNewRow,
            currentPage, pageSize, total, expandedRowKeys,
            selectedCompany, selectedCompanyId } = this.state;
        const { companiesData, countriesC2Data } = this.props
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
            type: 'radio',
        };

        let companiesWrappedData = []
        if (companiesData.data && companiesData.data.companies) {
            companiesWrappedData = [...companiesData.data.companies.map(item => { item.key = item.id; return item })]
        }
        if (hasNewRow) {
            companiesWrappedData = [{
                key: -1,
                id: -1,
                name: '',
                address: '',
                country_id: '',
            }, ...companiesWrappedData]
        } else {
            companiesWrappedData = [...companiesWrappedData.filter(item => item.id !== -1)]
        }

        let options = [];
        if (countriesC2Data.data && countriesC2Data.data.countries) {
            options = [...countriesC2Data.data.countries.map(item => { item.key = item.id; return item })]
        }

        const hasSelected = selectedRowKeys.length > 0 && selectedRowKeys[0] !== -1

        const companyColumns = [{
            title: '公司名',
            dataIndex: 'name',
            width: "25%",
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='company.name' value={record.name} onChange={this.onNewRowChange} />
                }
                return <a>{text}</a>
            }
        }, {
            title: '所在村',
            dataIndex: 'country_name',
            width: "15%",
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='company.country_id' value={record.country_id} onChange={this.onNewRowChange}
                        editType="select" valueType="int" options={options} placeholder="请选择村" />
                }
                return <a>{text}</a>
            }
        }, {
            title: '详细地址',
            dataIndex: 'address',
            width: "40%",
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell dataIndex='company.address' value={record.address} onChange={this.onNewRowChange} />
                }
                return <a href={record.url} target="_blank">{text}</a>
            }
        }, {
            title: '创建时间',
            dataIndex: 'create_at',
            width: "20%",
            render: (text, record) => {
                if (record.id === -1 || (editable && record.id === selectedRowKeys[0])) {
                    return <EditableCell type="opt" onSave={this.onRowSave} onCancel={this.handleCancelEditRow} />
                }
                var createAt = moment(new Date(text)).format(CONSTANTS.DATE_DISPLAY_FORMAT)
                return createAt;
            }
        }];


        return (
            <div className="gutter-example">
                <BreadcrumbCustom first="公司管理" />
                <CompanySearch fetchData={fetchData} />
                <Row gutter={16}>
                    <Col className="gutter-row" md={14}>
                        <div className="gutter-box">
                            <Card title="公司列表" bordered={false}>
                                <div style={{ marginBottom: 16 }}>
                                    <Button type="primary" onClick={this.handleAdd}
                                    >新增</Button>
                                    <Button type="primary" onClick={this.handleModify}
                                        disabled={!hasSelected}
                                    >修改</Button>
                                    <Button type="primary" onClick={this.showModal}
                                        disabled={!hasSelected}
                                    >删除</Button>
                                    <Modal
                                        title="警告"
                                        visible={this.state.visible}
                                        onOk={this.handleDelete}
                                        onCancel={this.hideModal}
                                        okText="确认"
                                        cancelText="取消"
                                    >
                                        <p>确认删除公司：{selectedCompany}</p>
                                    </Modal>
                                    <Button type="primary" onClick={this.downloadFile}
                                        disabled={loading}
                                    >下载</Button>
                                    <Upload style={{ marginLeft: '10px' }} {...this.uploadProps() }>
                                        <Button type="primary">上传
                                        </Button>
                                    </Upload>
                                </div>
                                <Table rowSelection={rowSelection}
                                    size="small"
                                    columns={companyColumns}
                                    dataSource={companiesWrappedData}
                                    onRow={(record) => ({
                                        onClick: () => this.onRowClick(record),
                                    })}
                                    pagination={{
                                        hideOnSinglePage: true,
                                        onChange: this.handlePageChange,
                                        current: currentPage,
                                        defaultCurrent: 1,
                                        pageSize,
                                        total,
                                    }}
                                />
                            </Card>
                        </div>
                    </Col>
                    <Col className="gutter-row" md={10}>
                        <div className="gutter-box">
                            <Card title={selectedCompany ? selectedCompany : "请选择公司"} bordered={false}
                                bodyStyle={{ paddingTop: 0 }}>
                                <div style={{}}>
                                </div>
                                {this.additionalTable()}
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
        companiesData = { data: { count: 0, towns: [] } },
        townsData = { data: { count: 0, towns: [] } },
        countriesData = { data: { count: 0, countries: [] } },
        countriesC2Data = { data: { count: 0, countries: [] } },
        usersInCompany = { data: { count: 0, users: [] } },
        placesInCompany = { data: { count: 0, monitor_places: [] } },
    } = state.httpData;
    return { companiesData, townsData, countriesC2Data, usersInCompany, placesInCompany };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(CompanyManager)

