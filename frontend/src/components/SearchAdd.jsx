import React, {useState} from 'react';
import {Button, Input, Modal, Popover} from 'antd';
import search from "../assets/images/search.png"
import add from "../assets/images/add.png"
import './SearchAdd.css'

const {TextArea} = Input;

const SearchAdd = ({searchAdd}) => {
    const {messageApi, sendMsg} = searchAdd

    const [show, setShowPop] = useState(false);
    const [dialogShow, setDialogShow] = useState(false);
    const [msgValue, setMsgValue] = useState("")
    const [uidValue, setUidValue] = useState("test_u2")

    const showPop = () => {
        if (show) {
            setShowPop(false);
        } else {
            setShowPop(true);
        }
    }

    const onPopItemClick = idx => {
        setShowPop(false)

        if (idx === 1) {
            messageApi.info("暂不支持")
            return
        }

        setDialogShow(true)
    }

    const popContent = (<div style={{
        display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center"
    }}>
        <span onClick={_ => onPopItemClick(0)} className="pop-item">发起单聊</span>
        <span onClick={_ => onPopItemClick(1)} className="pop-item">发起群聊</span>
    </div>)

    const onMsgSend = () => {
        const msd = {
            receiver: uidValue,
            chatType: 1,
            ttl: 0,
            msgContent: {
                type: 1,
                content: {
                    text: msgValue
                }
            }
        }

        sendMsg(msd)
    }

    return (<div>
        <div
            id="search-add"
            style={{
                width: '100%',
                height: 60,
                paddingLeft: 10,
                paddingRight: 10,
                backgroundColor: '#F7F7F7',
                boxSizing: "border-box",
                display: "flex",
                flexDirection: "row",
                alignItems: "center"
            }}>

            <div style={{
                // backgroundColor: 'red',
                boxSizing: "border-box",
                paddingRight: 8,
                display: "flex",
                flexDirection: "row",
                alignItems: "center",
                height: 28,
                flex: 1
            }}>
                <div style={{
                    display: "flex",
                    flexDirection: "row",
                    alignItems: "center",
                    backgroundColor: '#EAEAEA',
                    borderRadius: 3,
                    height: 28,
                    flex: 1
                }}>
                    <img src={search} style={{
                        width: 20, height: 20, marginLeft: 8, marginRight: 8,
                    }}/>

                    <span style={{
                        color: '#9E9E9E', fontSize: 13
                    }}>搜索</span>

                </div>

            </div>

            <Popover
                open={show}
                trigger="click"
                placement="bottomRight"
                content={popContent}
            >
                <div style={{
                    backgroundColor: "#EAEAEA",
                    width: 28,
                    height: 28,
                    alignItems: "center",
                    justifyContent: "center",
                    display: "flex",
                    borderRadius: 3,
                }} onClick={showPop}>
                    <img src={add} style={{
                        width: 16, height: 16,
                    }}/>
                </div>
            </Popover>

        </div>

        <Modal
            open={dialogShow}
            maskClosable={false}
            title={(<div style={{textAlign: "left"}}>发起单聊</div>)}
            okText="发送"
            cancelText="取消"
            closeIcon={false}
            footer={(<div style={{
                display: "flex", flexDirection: "row", justifyContent: "end", alignItems: "center", marginTop: 24,
            }}>
                <Button className="chat-modal-btn" style={{marginRight: 16}}
                        onClick={() => {
                            setDialogShow(false)
                        }}>取消</Button>
                <Button
                    disabled={!uidValue || !msgValue}
                    type="primary"
                    className="chat-modal-btn"
                    onClick={() => {
                        // setUidValue("")
                        setMsgValue("")

                        setDialogShow(false)

                        onMsgSend()
                    }}
                >发送</Button>
            </div>)}
        >
            <div style={{
                boxSizing: "border-box",
                marginTop: 24,
                display: "flex",
                flexDirection: "column",
                alignItems: "flex-start"
            }}>
                <Input
                    value={uidValue}
                    styles={{input: {width: 320, caretColor: '#1677FF'}}}
                    placeholder="Input Uid To Chat"
                    onChange={e => setUidValue(e.target.value)}
                />
                <TextArea
                    value={msgValue}
                    onChange={e => setMsgValue(e.target.value)}
                    placeholder="Input Message Content"
                    style={{
                        boxShadow: 'none',
                        marginTop: 16,
                        height: 120,
                        outline: 'none',
                        caretColor: '#1677FF',
                        resize: 'none'
                    }}
                />
            </div>
        </Modal>
    </div>)
}

export default SearchAdd