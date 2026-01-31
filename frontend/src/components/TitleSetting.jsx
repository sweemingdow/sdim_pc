import settingIcon from '../assets/images/setting.png'

const TitleSetting = ({setting, pointClick}) => {
    const {connState, title} = setting

    const getTitle = () => {
        if (connState === 1) {
            return {text: "未连接", color: "#666"}
        } else if (connState === 2) {
            return {text: "连接中...", color: "green"}
        } else if (connState === 3) {
            if (title) {
                return {text: title, color: "#161616"}
            }
            return {text: "", color: "#161616"}
        } else {
            return {text: "连接失败", color: "red"}
        }
    }

    const {text, color} = getTitle()

    return (<div style={{
        width: '100%',
        backgroundColor: "#EDEDED",
        height: 59,
        display: "flex",
        flexDirection: "row",
        alignItems: "center",
        boxSizing: "border-box",
        paddingLeft: 20,
        paddingRight: 20
    }}>


        <span style={{
            color: color,
            fontSize: 16,
            fontWeight: "bold",
        }}>
            {text}
        </span>

        <div style={{flex: 1}}/>

        <div style={{
            width: 28,
            height: 28,
            backgroundColor: "#E1E1E1",
            borderRadius: 3,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
        }}
             onClick={pointClick}>
            <img src={settingIcon} style={{
                width: 18,
                height: 18,
            }}/>
        </div>

    </div>)
}
export default TitleSetting