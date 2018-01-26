/**
 * Created by hao.cheng on 2017/4/21.
 */
import React, { Component } from 'react';
import ReactEcharts from 'echarts-for-react';


class PieChart extends Component {

    genOption = chartOptions => {
        if(chartOptions === undefined)return
        let legendData = []
        for(let option of chartOptions){
            legendData.push(option.name)
        }
        
        return { 
            title : {
                //text: '某站点用户访问来源',
                //subtext: '纯属虚构',
                x:'center'
            },
            tooltip : {
                trigger: 'item',
                formatter: "{a} <br/>{b} : {c} ({d}%)"
            },
            legend: {
                orient: 'vertical',
                left: 'left',
                data: legendData
            },
            series : [
                {
                    name: ' 完成情况',
                    type: 'pie',
                    radius : '55%',
                    center: ['50%', '60%'],
                    data: chartOptions,
                    itemStyle: {
                        emphasis: {
                            shadowBlur: 10,
                            shadowOffsetX: 0,
                            shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                }
            ]
        }
    }

    render(){
        const { chartOptions } = this.props
        return (
            <ReactEcharts
                theme="light"
                option={this.genOption(chartOptions)}
                style={{ width: '100%'}}
            />
        )
    }
}

export default PieChart;