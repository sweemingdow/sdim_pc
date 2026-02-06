import GroupSideCurtain from "./GroupSideCurtain.jsx"
import P2pSideCurtain from "./P2pSideCurtain.jsx"

const SideCurtain = ({convItem, messageApi}) => {
    // console.dir(convItem)
    const convType = convItem?.convType
    if (convType === 2) {
        return (<GroupSideCurtain convItem={convItem}
                                  messageApi={messageApi}/>)
    } else if (convType === 1) {
        return (<P2pSideCurtain/>)
    }
}

export default SideCurtain