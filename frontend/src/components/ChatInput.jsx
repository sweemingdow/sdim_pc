import React, {useState} from 'react';

import {Button, Input} from 'antd';

const {TextArea} = Input;

const ChatInput = ({send, inputEnabled, onClick}) => {
    const sendMsg = send
    const [value, setInputValue] = useState('');

    const onInputValueChanged = (e) => {
        let val = e.target.value

        setInputValue(val)
    }

    const isValBlank = val => {
        return !val || String(val).trim() === ''
    }


    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            if (isValBlank(value)) {
                if (e.shiftKey) {
                    // Shift + Enter：允许默认行为（换行）
                    return;
                }

                e.preventDefault();
                return;
            }

            // 非空, 允许shift+enter 换行
            if (e.shiftKey) {
                return;
            }

            e.preventDefault();
            onSendMsg();
        }
    };

    const onSendMsg = (_) => {
        const msd = {
            chatType: 1,
            ttl: 0,
            msgContent: {
                type: 1,
                content: {
                    text: value
                }
            }
        }

        setInputValue("")
        sendMsg(msd)
    }

    return (<div style={{
        width: '100%',
        height: 180,
        backgroundColor: "cyan",
        display: "flex",
        flexDirection: "column",
    }}>

        <div style={{
            height: 50,
            display: "none",
            flexDirection: "column",
            alignItems: "center",
            backgroundColor: "white"
        }}></div>

        <div style={{
            paddingTop: 8,
            paddingLeft: 5,
            paddingRight: 5,
            paddingBottom: 12,
            flex: 1,
            backgroundColor: "#EDEDED",
            display: "flex",
            flexDirection: "column",
        }}>

            <TextArea
                disabled={!inputEnabled}
                bordered={false}
                style={{
                    flex: 1,
                    border: 0,
                    boxShadow: 'none',
                    outline: 'none',
                    caretColor: '#1677FF',
                    resize: 'none'
                }}
                value={value}
                onKeyDown={handleKeyDown}
                onChange={onInputValueChanged}
                onClick={onClick}
            />

            <Button
                style={{
                    alignSelf: "flex-end",
                    width: 90,
                    height: 32,
                    borderRadius: 4,
                }}
                disabled={isValBlank(value)}
                onClick={onSendMsg}
                type="primary">发送(S)</Button>

        </div>


    </div>)
}

export default ChatInput