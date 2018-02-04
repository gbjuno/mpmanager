/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card, Pagination, Icon, Tabs, Modal, message } from 'antd';
import moment from 'moment';
import * as _ from 'lodash';
import { fetchData, receiveData } from '../../action';
import * as CONSTANTS from '../../constants';
import BreadcrumbCustom from '../../components/BreadcrumbCustom';
import PlaceSearch from '../search/PlaceSearch'
import PlaceForm from '../forms/PlaceForm'
import * as config from '../../axios/config'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

const TabPane = Tabs.TabPane;

class PlaceManager extends React.Component {
    state = {
        placeTypes: [],
        placesData: [],
        placesDataWithType: [],
        rate: 1,
        placeTypeLoading: false,
        standardHeight: 276,
        hasNew: false,
        deleteModal: false,
        deleteRecord: {},
        selectedCompanyId: '',
        selectedCompanyName: '',
        newPlaceName: '',
        selectedPlaceTypeId: '',
    };

    componentDidMount = () => {
        this.resizeWindow();
        window.onresize = () =>{
            this.resizeWindow();
        };

        this.fetchPlaceType();
    };

    fetchPlaceType = () => {
        const { fetchData } = this.props
        let tempTownId
        fetchData({funcName: 'fetchPlaceTypes', stateName: 'placeTypes'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.monitor_types === undefined) return
            this.setState({
                placeTypes: [...res.data.monitor_types.map(val => {
                    val.key = val.id;
                    return val;
                })],
                placeTypeLoading: false,
                selectedPlaceTypeId: res.data.monitor_types[0].id,
            }, () => {
                //this.fetchPlaceData();
            });
        });
    }

    fetchPlaceData = () => {
        const { fetchData } = this.props
        const { placeTypeLoading, placeTypes } = this.state
        fetchData({funcName: 'fetchPlaces', stateName: 'placesData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.monitor_places === undefined) return
            let placesDataWithType = []
            for(let placeType of placeTypes){
                placesDataWithType.push(
                    {
                    placeTypeId: placeType.id,
                    placeTypeName: placeType.name,
                    placesData: [...res.data.monitor_places.map(val => {
                            val.key = val.id;
                            return val;
                        }).filter(val => val.monitor_type_id === placeType.id)],
                    }
                );
            }
            this.setState({
                placesDataWithType,
            });
        });
    }

    searchPlace = (companyId) => {
        const { fetchData, onSearch } = this.props
        fetchData({funcName: 'searchPlaces', params: {companyId}, 
            stateName: 'placesData'})
            .then(res => {
                if(onSearch){
                    onSearch()
                }
            })
    }

    handleAdd = () => {
        this.setState({
            hasNew: true,
        })
    }

    handleFormChange = (name) => {
        this.setState({
            newPlaceName: name,
        })
    }

    handleTabChange = (placeTypeId) => {
        this.setState({
            hasNew: false,
            selectedPlaceTypeId: placeTypeId,
        })
    }

    handleSearchChange = (companyId, companyName) => {
        this.setState({
            selectedCompanyId: companyId,
            selectedCompanyName: companyName,
        })
    }

    handleSearchSubmit = () => {
        this.setState({
            hasNew: false,
        })
    }

    handleSave = () => {
        const { fetchData } = this.props
        const { newPlaceName, selectedCompanyId, selectedPlaceTypeId } = this.state
        let saveObj = {
            name: newPlaceName,
            company_id: parseInt(selectedCompanyId),
            monitor_type_id: parseInt(selectedPlaceTypeId),
        }
        fetchData({funcName: 'newPlace', params: saveObj, 
            stateName: 'newPlaceStatus'}).then(res => {
                message.success('创建成功')
                this.setState({
                    hasNew: false
                })
                this.searchPlace(selectedCompanyId)
            }).catch(err => {
                let errRes = err.response
                if(errRes && errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            }
        );
        
        //this.fetchPlaceData()
    }

    handleCancelEdit = () => {
        this.setState({
            hasNew: false,
        })
    }

    handleDelete = (record, e) => {
        this.setState({
            deleteModal: true,
            deleteRecord: record,
        })
    }

    handleDeleteConfirm = () => {
        const { deleteRecord, selectedCompanyId } = this.state
        const { fetchData }  = this.props
        fetchData({
            funcName: 'deletePlace', params: { id: deleteRecord.id }, stateName: 'deletePlaceStatus'
            }).then(res => {
                message.success('删除成功')
                this.searchPlace(selectedCompanyId)
                this.setState({
                    deleteModal: false,
                    deleteRecord: {},
                })
            }).catch(err => {
                let errRes = err.response
                if(errRes && errRes.data && errRes.data.status === 'error'){
                    message.error(errRes.data.error)
                }
            });
    }

    hideDeleteModal = () => {
        this.setState({
            deleteModal: false,
            deleteRecord: {},
        })
    }

    componentDidUpdate = (nextProps, nextState) => {
        
    };

    
    getClientWidth = () => {    // 获取当前浏览器宽度并设置responsive管理响应式
        const { receiveData } = this.props;
        const clientWidth = document.body.clientWidth;
        receiveData({isMobile: clientWidth <= 992}, 'responsive');
    };


    resizeWindow = () => {
        this.getClientWidth();
        const placeQRs = document.getElementById("placeQRs");
        if(placeQRs === undefined || placeQRs === null) return;
        const swidth = document.body.clientWidth - 200;
        const benchmark = 1680
        this.setState({
            contentWidth: swidth,
            rate: swidth / benchmark,
        });
        
    }

    //转置图片数据
    transpositionToMatrix = picArray => {
        if(picArray===undefined || picArray.length === undefined || picArray.length === 0) return [[]];
        const colLen = 6
        const rowLen = Math.ceil(picArray.length / colLen)
        const mod = picArray.length % colLen
        let matrix = new Array()
        for(let i=0; i<colLen; i++){
            let uniArray = new Array()
            for(let j=0; j<rowLen; j++)
            {
                if(j * colLen + i >= picArray.length) break;
                uniArray.push(picArray[j * colLen + i])
            }
            matrix.push(uniArray)
        }
        return matrix
    };


    generateCard = (imgs, isMobile=false) => imgs.map(v1 => (
        v1.map(v2 => { 
            const { standardHeight, rate, hasNew } = this.state
            let baseHeight = standardHeight * rate
            return (
            <div key={v2.id} className="gutter-box" style={isMobile? {}: {height: baseHeight + 120}}>
                {v2.id === -1 ?
                    hasNew?
                    <Card bordered={false} style={{ }} 
                    bodyStyle={isMobile? {}: {cursor: 'pointer', verticalAlign: 'middle', height: baseHeight + 60 }}
                    actions={[<a onClick={this.handleSave}>保存</a>, <a onClick={this.handleCancelEdit}>取消</a>]}>
                        <PlaceForm  onChange={this.handleFormChange} value={v2}/>
                    </Card>
                    :
                    <Card bordered={false} style={{border: '1px dashed #d9d9d9', textAlign: 'center'}} onClick={this.handleAdd}
                    bodyStyle={isMobile? {}: {paddingTop: 123, cursor: 'pointer', verticalAlign: 'middle', height: baseHeight + 60 + 45}}
                    >
                        <Icon style={{fontSize: 33}} type="plus"/>
                    </Card>
                :
                <Card bordered={false} bodyStyle={isMobile? {padding: 0}: { padding: 0, height: baseHeight + 60}}
                actions={[<Icon type="edit" />, <Icon type="delete" onClick={this.handleDelete.bind(this, v2)} />]}>
                    <div>
                        <img style={isMobile? {}: {height: baseHeight}} onClick={() => {}} 
                            alt="example" width="100%" src={config.SERVER_ROOT + v2.qrcode_uri} />
                    </div>
                    <div className="pa-s">
                        <small><a>{v2.name}<span style={{paddingLeft: 5}}>{moment(new Date(v2.create_at)).format(CONSTANTS.DATE_DISPLAY_FORMAT)}</span></a></small>
                    </div>
                </Card>
                }
            </div>
        )
        })
    ))

    generateGrid = (datasWithType=[], isMobile) => datasWithType.map(dataWithType => {
        let imgs = this.transpositionToMatrix( dataWithType.placesData);
        const imgsTag = this.generateCard(imgs, isMobile)
        return (
        <TabPane tab={dataWithType.placeTypeName} key={dataWithType.placeTypeId} >
            <Row gutter={20}>
                <Col className="gutter-row" md={4}>
                    {imgsTag[0]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[1]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[2]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[3]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[4]}
                </Col>
                <Col className="gutter-row" md={4}>
                    {imgsTag[5]}
                </Col>
            </Row>
        </TabPane>
        )
    })

    transform = (placesData, placeTypes) => {
        const {selectedCompanyId, selectedCompanyName} = this.state
        if(placeTypes === undefined || placesData === undefined 
            || _.isEmpty(placeTypes) || _.isEmpty(placesData)
            || placesData.data === undefined || placesData.data.monitor_places === undefined) {
            return []
        }else{
            let placesDataWithType =[]
            let places = placesData.data
            for(let placeType of placeTypes){
                placesDataWithType.push({
                    placeTypeId: placeType.id,
                    placeTypeName: placeType.name,
                    placesData: [...places.monitor_places.map(val => {
                        val.key = val.id;
                        return val;
                    }).filter(val => val.monitor_type_id === placeType.id), {
                        id: -1,
                        company_id: selectedCompanyId,
                        company_name: selectedCompanyName,
                        monitor_type_id: placeType.id,
                        monitor_type_name: placeType.name,
                        qrcode_path: null,
                        qrcode_uri: null,
                    }],
                });
            }
            return placesDataWithType
        }
    }

    render() {
        const { rate, placeTypes } = this.state
        const { placesData, responsive } = this.props
        let isMobile = responsive.data.isMobile
        
        let placesDataWithType = this.transform(placesData, placeTypes)
        let placeGrids = this.generateGrid(placesDataWithType, isMobile)
        let defaultTabId = placeTypes[0]?`${placeTypes[0].id}`:'0'
        
        return (
            <div id="placeQRs" className="gutter-example button-demo">
                <BreadcrumbCustom first="地点管理" second="" />
                <PlaceSearch onChange={this.handleSearchChange} onSearch={this.handleSearchSubmit} fetchData={fetchData}/>
                <Tabs defaultActiveKey={defaultTabId} 
                    onChange={this.handleTabChange}>
                {placeGrids}
                </Tabs>
                <Modal
                    title="警告"
                    visible={this.state.deleteModal}
                    onOk={this.handleDeleteConfirm}
                    onCancel={this.hideDeleteModal}
                    okText="确认"
                    cancelText="取消"
                    >
                    <p>删除地点：{this.state.deleteRecord.name}?</p>
                </Modal>
                <style>{`
                    .ant-card-body img {
                        cursor: pointer;
                    }
                `}</style>
            </div>
        )
    }
}

const mapStateToProps = state => {
    const { placesData = {data: []} } = state.httpData;
    return { ...state.httpData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PlaceManager);