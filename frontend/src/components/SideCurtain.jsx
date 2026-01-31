import GroupSideCurtain from "./GroupSideCurtain.jsx"
import P2pSideCurtain from "./P2pSideCurtain.jsx"

const SideCurtain = ({convItem}) => {
    // console.dir(convItem)
    const convType = convItem?.convType
    if (convType === 2) {
        return (<GroupSideCurtain convItem={convItem}/>)
    } else if (convType === 1) {
        return (<P2pSideCurtain/>)
    }
}

export default SideCurtain