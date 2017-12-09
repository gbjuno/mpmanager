/**
 * Created by Jingle on 2017/11/4.
 */
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Row, Col, Card } from 'antd';
import { fetchData, receiveData } from '../../action';
import BreadcrumbCustom from '../BreadcrumbCustom';
import PlaceSearch from './search/PlaceSearch'
import * as config from '../../axios/config'
import PhotoSwipe from 'photoswipe';
import PhotoswipeUIDefault from 'photoswipe/dist/photoswipe-ui-default';

import 'photoswipe/dist/photoswipe.css';
import 'photoswipe/dist/default-skin/default-skin.css';

class PlaceManager extends React.Component {
    state = {
        placeTypes: [],
        placesData: [],
        placesDataWithType: [],
        rate: 1,
        responsive: false,
        placeTypeLoading: false,
        standardHeight: 300
    };

    componentWillUnmount = () => {
        
    };

    componentDidMount = () => {
        this.setState({ placeTypeLoading: true });
        this.resizePicture();
        const clientWidth = document.body.clientWidth;
        if(clientWidth <= 992) {
            this.setState({
                responsive: true,
            })
        }
        window.onresize = () =>{
            const clientWidth = document.body.clientWidth;
            if(clientWidth <= 992) {
                this.setState({
                    responsive: true,
                })
                return;
            }else{
                this.resizePicture();
            }
            
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
            }, () => {
                this.fetchPlaceData();
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
                console.log('placeType.id', placeType.id)
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

    componentDidUpdate = (nextProps, nextState) => {
    };

    


    resizePicture = () => {
        const placeQRs = document.getElementById("placeQRs");
        if(placeQRs === undefined || placeQRs === null) return;
        const swidth = placeQRs.clientWidth;
        const benchmark = 1680;
        this.setState({
            rate: swidth / benchmark,
            responsive: false,
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

    generateCard = imgs => imgs.map(v1 => (
        v1.map(v2 => (
            <div key={v2.id} className="gutter-box" style={this.state.responsive? {}: {height: this.state.standardHeight * this.state.rate + 80}}>
                <Card bordered={false} bodyStyle={this.state.responsive? {padding: 0}: { padding: 0, height: this.state.standardHeight * this.state.rate + 60}}>
                    <div>
                        <img style={this.state.responsive? {}: {height: this.state.standardHeight * this.state.rate}} onClick={() => {}} 
                            alt="example" width="100%" src={config.SERVER_ROOT + v2.qrcode_uri} />
                    </div>
                    <div className="pa-m">
                        <h3>{v2.companyName}<span style={{paddingLeft: 5}}>{v2.name}</span></h3>
                        <small><a>{v2.placeName}<span style={{paddingLeft: 5}}>{v2.create_at}</span></a></small>
                    </div>
                </Card>
            </div>
            ))
        ))

    render() {
        const { rate, responsive, placeTypes, placesData, placesDataWithType } = this.state
        console.log('placesDataWithType...', placesDataWithType)
        console.log('placeTypes...', placeTypes)
        let placeGrids = placesDataWithType.map(placeDataWithType => {
            let imgs = this.transpositionToMatrix( placeDataWithType.placesData);
            
            const imgsTag = this.generateCard(imgs)
            return (
            <div>
                <div>{placeDataWithType.name}</div>
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
            </div>
            )
        })
        
        return (
            <div id="placeQRs" className="gutter-example button-demo">
                <BreadcrumbCustom first="安监管理" second="地点管理" />
                <PlaceSearch style={{paddingBottom: 13}} fetchData={fetchData}/>
                {placeGrids}
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
    const { picData = {data: []} } = state.httpData;
    return { picData };
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(PlaceManager);