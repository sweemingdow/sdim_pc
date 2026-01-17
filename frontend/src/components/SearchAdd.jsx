import React, {useState} from 'react';
import {Button, Input, Modal, Popover} from 'antd';
import search from "../assets/images/search.png"
import add from "../assets/images/add.png"
import './SearchAdd.css'

const {TextArea} = Input;

const SearchAdd = ({searchAdd}) => {
    const {messageApi, sendMsg, startGroupChat} = searchAdd

    const [show, setShowPop] = useState(false);
    const [dialogShow, setDialogShow] = useState(false);
    const [msgValue, setMsgValue] = useState("")
    const [uidValue, setUidValue] = useState("test_u2")
    const [p2pChat, setP2pChat] = useState(true)
    const [groupName, setGroupName] = useState("grp_test")
    const [groupMembers, setGroupMembers] = useState("test_u2,test_u3")
    const [limitNum, setLimitNum] = useState("2")
    const [groupCreating, setGroupCreating] = useState(false)

    const showPop = () => {
        if (show) {
            setShowPop(false);
        } else {
            setShowPop(true);
        }
    }

    const onPopItemClick = idx => {
        setShowPop(false)
        setP2pChat(idx === 0);

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
            receiver: uidValue, chatType: 1, ttl: 0, msgContent: {
                type: 1, content: {
                    text: msgValue
                }
            }
        }

        sendMsg(msd)
    }


    const foot = () => {
        if (p2pChat) {
            return (<div style={{
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
            </div>)
        } else {
            return (<div style={{
                display: "flex", flexDirection: "row", justifyContent: "end", alignItems: "center", marginTop: 24,
            }}>
                <Button className="chat-modal-btn" style={{marginRight: 16}}
                        onClick={() => {
                            setDialogShow(false)
                        }}>取消</Button>
                <Button
                    disabled={!groupName || !groupMembers || !limitNum || groupCreating}
                    type="primary"
                    loading={groupCreating}
                    className="chat-modal-btn"
                    style={{width: 100}}
                    onClick={async () => {
                        // type StartGroupChatData struct {
                        // 	GroupName  string `json:"groupName"`
                        // 	Avatar     string `json:"avatar"`
                        // 	LimitedNum string    `json:"limitedNum"`
                        // 	MembersStr string `json:"membersStr"`
                        // }
                        const data = {
                            groupName,
                            limitedNum: limitNum,
                            membersStr: groupMembers
                        };

                        setGroupCreating(true)

                        await startGroupChat(data)

                        // setGroupName("")
                        // setGroupMembers("")
                        // setLimitNum("")
                        setGroupCreating(false)

                        setDialogShow(false)
                    }}
                >创建群聊</Button>
            </div>)
        }
    }

    const modalBody = () => {
        if (p2pChat) {
            return (<div style={{
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
            </div>)
        } else {
            return (<div style={{
                boxSizing: "border-box",
                marginTop: 24,
                display: "flex",
                flexDirection: "column",
                alignItems: "flex-start"
            }}>
                <Input
                    value={groupName}
                    styles={{input: {width: 320, caretColor: '#1677FF'}}}
                    placeholder="Input GroupName To Create"
                    onChange={e => setGroupName(e.target.value)}
                />
                <Input
                    value={limitNum}
                    styles={{input: {width: 120, caretColor: '#1677FF', marginTop: 16,}}}
                    placeholder="Limited Num"
                    onChange={e => setLimitNum(e.target.value)}
                />
                <TextArea
                    value={groupMembers}
                    onChange={e => setGroupMembers(e.target.value)}
                    placeholder="Input Users For Create Group(Like: u1;u2;u3)"
                    style={{
                        boxShadow: 'none',
                        marginTop: 16,
                        height: 120,
                        outline: 'none',
                        caretColor: '#1677FF',
                        resize: 'none'
                    }}
                />
            </div>)
        }
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
            title={(<div style={{textAlign: "left"}}>{p2pChat ? '发起单聊' : '发起群聊'}</div>)}
            closeIcon={false}
            footer={foot()}
        >
            {modalBody()}
        </Modal>
    </div>)
}

export default SearchAdd