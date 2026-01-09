import './App.css';

import ConvList from "./components/ConvList.jsx";
import Menu from "./components/Menu.jsx";
import SearchAdd from "./components/SearchAdd.jsx";
import TitleSetting from "./components/TitleSetting.jsx";
import ChatRoom from "./components/ChatRoom.jsx";
import ChatInput from "./components/ChatInput.jsx";
import DebugRowOne from "./components/DebugRowOne.jsx";
import {message} from 'antd'
import {Conn2Engine, Disconnect, SendMsg} from "../wailsjs/go/main/App.js";
import {useRef, useState} from "react";
import emptyIcon from './assets/images/room_empty.png'
import {SyncConvList} from "../wailsjs/go/syncbinder/SyncBinder.js";

function App() {
    const [messageApi, contextHolder] = message.useMessage();
    const convListRef = useRef(null);

    /*
    * connState:
    * 1=未连接
    * 2=连接中
    * 3=连接成功
    * 4=连接失败
    * */
    const [connState, setConnState] = useState(1)

    const [maskShow, setMaskShow] = useState(false)

    const [selectedConvItem, setSelectedConvItem] = useState({})

    // 主动连接
    const connWithActive = async (uid) => {
        setConnState(2)

        try {
            // 建立连接
            await Conn2Engine(uid)

            setConnState(3)

            messageApi.info(`uid=${uid}, 连接成功`)

            // 拉取会话列表
            let convItems = await SyncConvList(uid)

            if (convListRef.current && convItems) {
                console.log(`sync conv list success, convItems:${JSON.stringify(convItems)}`)
                
                convListRef.current.setConvItemsExternal(convItems)
            }

        } catch (e) {
            setConnState(4)
            messageApi.error(`连接错误, uid=${uid}, e=${e}`)
        }

    }

    // 主动断开连接
    const disConnWithActive = () => {
        setConnState(1)

        Disconnect().then(_ => {
            messageApi.info(`断连成功`)
        }).catch(e => {
            messageApi.error(`断连失败, e=${e}`)
        })
    }


    // 发送消息
    // type MsgSendData struct {
    //     //Sender     string           `json:"sender,omitempty"`   // 发送者uid
    //     Receiver   string           `json:"receiver,omitempty"` // 接收者, 单聊是对方的uid, 群聊是群id
    //     ChatType   preinld.ChatType `json:"chatType,omitempty"`
    //     Ttl        int32            `json:"ttl,omitempty"` // 消息过期时间(sec), -1:阅后即焚,0:不过期
    //     MsgContent *msg.MsgContent  `json:"msgContent,omitempty"`
    // }
    const sendMsg = (msd) => {
        SendMsg(msd).then(_ => {
            messageApi.success(`msg send success, msd=${JSON.stringify(msd)}`)
        }).catch(e => {
            messageApi.error(`msg send failed, msg=${msd}, e=${e}`)
        })
    }

    const onConvItemSelected = convItem => {
        console.log(`on conv item selected, convItem:${JSON.stringify(convItem)}`)
        setSelectedConvItem(convItem)

        // todo 同步该会话中的消息
    }

    return (<>
        {contextHolder}
        <div id="App">
            <div id="outer" style={{
                width: '100vw', height: '100vh', // backgroundColor: 'cyan',
                display: "flex", flexDirection: "column",
            }}>
                <div id="dev-top" style={{
                    width: '100%', flex: 1, display: "flex", flexDirection: "column", justifyContent: "flex-start"
                }}>
                    <DebugRowOne rowOne={{connState, messageApi, connWithActive, disConnWithActive}}/>

                    <div style={{
                        height: 1, width: "100%", backgroundColor: '#ddd'
                    }}/>

                </div>

                <div id="bottom" style={{
                    flex: 9, display: "flex", flexDirection: "row",
                }}>
                    <div id="b-l" style={{
                        backgroundColor: "#eee", height: '100%', boxSizing: "border-box", padding: 10,
                    }}>
                        <Menu></Menu>
                    </div>

                    <div style={{
                        height: '100%', width: 1, backgroundColor: '#ccc'
                    }}>

                    </div>
                    <div id="b-m" style={{
                        flex: 2.2, minWidth: 0, backgroundColor: '#F7F7F7', overflow: "hidden",
                    }}>
                        <SearchAdd searchAdd={{messageApi, sendMsg}}/>

                        <ConvList
                            ref={convListRef}
                            setMaskShow={setMaskShow}
                            setSelectedConvItem={onConvItemSelected}/>
                    </div>

                    <div style={{
                        height: '100%', width: 1, backgroundColor: '#ddd'
                    }}>
                    </div>
                    <div id="b-r" style={{
                        flex: 7.3, display: "flex", flexDirection: "column", position: "relative"
                    }}>

                        <div style={{
                            position: "absolute",
                            width: '100%',
                            height: '100%',
                            backgroundColor: "#EDEDED",
                            zIndex: 999,
                            display: maskShow ? "flex" : "none",
                            justifyContent: "center",
                            alignItems: "center"
                        }}>
                            <img src={emptyIcon}
                                 style={{
                                     width: 64,
                                     height: 64,
                                 }}
                            />
                        </div>
                        <TitleSetting setting={{connState, title: selectedConvItem ? selectedConvItem.title : ""}}/>

                        <div style={{
                            height: 1, width: "100%", backgroundColor: '#ddd'
                        }}/>

                        <div style={{
                            backgroundColor: "red", flex: 1, display: "flex", flexDirection: "column"
                        }}>

                            <ChatRoom/>

                            <div style={{
                                height: 1, width: "100%", backgroundColor: '#ddd'
                            }}/>

                            <ChatInput send={{sendMsg}}/>

                        </div>

                    </div>

                </div>
            </div>
        </div>
    </>)
}

export default App