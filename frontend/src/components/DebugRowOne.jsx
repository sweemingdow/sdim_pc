import React, {useState} from 'react';
import {Button, Input} from 'antd';

const DebugRowOne = ({rowOne}) => {
    /*
    * connState:
    * 1=未连接
    * 2=连接中
    * 3=连接成功
    * 4=连接失败
    * */
    const {connState, messageApi, connWithActive, disConnWithActive} = rowOne

    const [uidValue, setUidValue] = useState("test_u1")

    const onConnectClick = (e) => {
        if (!uidValue) {
            messageApi.warning("未输入Uid")
            return
        }

        connWithActive(uidValue)
    }

    const onDisConnClick = (e) => {
        disConnWithActive()
    }


    return (<div style={{
        width: '100%',
        height: 34,
        display: "flex",
        boxSizing: "border-box",
        paddingLeft: 16,
        paddingRight: 16,
        flexDirection: "row",
        alignItems: "center"
    }}>

        <Input
            value={uidValue}
            onChange={e => setUidValue(e.target.value)}
            style={{
                width: 180
            }}
            size="small"
            placeholder="Input Uid"
        />
        <Button
            disabled={connState === 2 || connState === 3}
            onClick={onConnectClick}
            loading={connState === 2}
            style={{
                borderRadius: 3,
                marginLeft: 12,
                height: 24,
            }}>Connect</Button>

        <Button
            disabled={connState !== 3}
            onClick={onDisConnClick}
            style={{
                borderRadius: 3,
                marginLeft: 24,
                height: 24,
                width: 80,
            }}>Disconnect</Button>
    </div>)
}

export default DebugRowOne