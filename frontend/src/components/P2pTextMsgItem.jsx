import {LoadingOutlined} from '@ant-design/icons';
import {Spin} from 'antd';
import '../css/msg_item_time.css'
import failedIcon from '../assets/images/sent_failed.png'

const P2pTextMsgItem = ({msg, pw, _}) => {
    const isSelf = msg.isSelf
    const maxWidth = pw * 0.65

    const timeEleStyle = () => {
        return {display: msg.timeShow ? 'block' : 'none'}
    }

    // console.log(`msg=${JSON.stringify(msg)}`)

    const content = msg.content.content
    const sendInfo = msg.senderInfo
    const text = content.text

    if (isSelf) {
        return (<div className="msg-item-outer">
            <div
                style={timeEleStyle()}
                className="msg-item-time">
                {msg.timeText}
            </div>
            <div style={{
                alignSelf: "flex-end",
                boxSizing: "border-box",
                marginRight: 18,
                display: "flex",
                flexDirection: "row-reverse",
                marginBottom: 8,
                minHeight: 34,
                maxWidth: maxWidth,
                justifyContent: "reverse",
            }}>
                <img
                    src={sendInfo.avatar}
                    style={{
                        flexShrink: 0, width: 34, height: 34, marginLeft: 10,
                    }}>
                </img>

                <div style={{
                    width: 0,
                    height: 0,
                    borderTop: '6px solid transparent',
                    borderBottom: '6px solid transparent',
                    borderLeft: '6px solid #8ce99a',
                    marginTop: 11,
                }}></div>


                <div style={{
                    boxSizing: "border-box",
                    display: "flex",
                    flex: 1,
                    borderRadius: 5,
                    backgroundColor: "#8ce99a",
                    paddingTop: 8,
                    paddingBottom: 8,
                    paddingLeft: 10,
                    paddingRight: 10,
                    minHeight: 34,
                    textAlign: "start",
                    border: 'none',
                    fontSize: 14,
                    color: "#161616",
                    wordBreak: 'break-word',
                    whiteSpace: 'pre-wrap',
                }}>{text}</div>


                <div style={{
                    alignSelf: "center",
                    marginRight: 10,
                    display: msg.state === 2 ? "none" : "flex",
                    justifyContent: "center",
                    alignItems: "center"
                }}>
                    {
                        msg.state === 1 ? <Spin
                                spinning={true}
                                delay={350}
                                indicator={<LoadingOutlined spin/>} size="small"/> :
                            <img src={failedIcon}
                                 style={{
                                     width: 18,
                                     height: 18
                                 }}/>
                    }


                </div>

            </div>
        </div>);
    } else {
        return (<div className="msg-item-outer">
            <div
                style={timeEleStyle()}
                className="msg-item-time">
                {msg.timeText}
            </div>
            <div style={{
                alignSelf: "flex-start",
                boxSizing: "border-box",
                marginRight: 18,
                display: "flex",
                flexDirection: "row",
                marginBottom: 8,
                minHeight: 34,
                maxWidth: maxWidth,
            }}>
                <img
                    src={sendInfo.avatar}
                    style={{
                        flexShrink: 0, width: 34, height: 34, marginLeft: 18,
                    }}>
                </img>

                <div style={{
                    width: 0,
                    height: 0,
                    borderTop: '6px solid transparent',
                    borderBottom: '6px solid transparent',
                    borderRight: '6px solid white',
                    marginLeft: 10,
                    marginTop: 11,
                }}></div>


                <div style={{
                    boxSizing: "border-box",
                    display: "flex",
                    alignItems: "center",
                    paddingTop: 8,
                    flex: 1,
                    paddingBottom: 8,
                    paddingLeft: 10,
                    paddingRight: 10,
                    borderRadius: 5,
                    backgroundColor: "white",
                    minHeight: 34,
                    textAlign: "start",
                    border: 'none',
                    fontSize: 14,
                    color: "#161616",
                    wordBreak: 'break-word',
                    whiteSpace: 'pre-wrap',
                }}>{text}</div>

                <div style={{
                    alignSelf: "center",
                    marginLeft: 10,
                    display: msg.state === 2 ? "none" : "flex",
                    justifyContent: "center",
                    alignItems: "center"
                }}>
                    {
                        msg.state === 1 ? <Spin
                                spinning={true}
                                delay={350}
                                indicator={<LoadingOutlined spin/>} size="small"/> :
                            <img src={failedIcon}
                                 style={{
                                     width: 18,
                                     height: 18
                                 }}/>
                    }


                </div>
            </div>
        </div>);
    }
}

export default P2pTextMsgItem