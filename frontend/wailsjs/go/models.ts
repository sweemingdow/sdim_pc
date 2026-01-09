export namespace conv {
	
	export class ConvItem {
	    convId?: string;
	    convType?: number;
	    icon?: string;
	    title?: string;
	    relationId?: string;
	    remark?: string;
	    pinTop?: boolean;
	    noDisturb?: boolean;
	    msgSeq?: number;
	    lastMsg?: preinld.Msg;
	    browseMsgSeq?: number;
	    unreadCount?: number;
	    cts?: number;
	    uts?: number;
	    recentlyMsgs: preinld.Msg[];
	
	    static createFrom(source: any = {}) {
	        return new ConvItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.convId = source["convId"];
	        this.convType = source["convType"];
	        this.icon = source["icon"];
	        this.title = source["title"];
	        this.relationId = source["relationId"];
	        this.remark = source["remark"];
	        this.pinTop = source["pinTop"];
	        this.noDisturb = source["noDisturb"];
	        this.msgSeq = source["msgSeq"];
	        this.lastMsg = this.convertValues(source["lastMsg"], preinld.Msg);
	        this.browseMsgSeq = source["browseMsgSeq"];
	        this.unreadCount = source["unreadCount"];
	        this.cts = source["cts"];
	        this.uts = source["uts"];
	        this.recentlyMsgs = this.convertValues(source["recentlyMsgs"], preinld.Msg);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace preinld {
	
	export class SenderInfo {
	    nickname?: string;
	    avatar?: string;
	
	    static createFrom(source: any = {}) {
	        return new SenderInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nickname = source["nickname"];
	        this.avatar = source["avatar"];
	    }
	}
	export class MsgContent {
	    type?: number;
	    content?: Record<string, any>;
	    custom?: Record<string, any>;
	    extra?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new MsgContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.content = source["content"];
	        this.custom = source["custom"];
	        this.extra = source["extra"];
	    }
	}
	export class Msg {
	    msgId: number;
	    convId: string;
	    sender: string;
	    receiver: string;
	    chatType: number;
	    msgType: number;
	    content?: MsgContent;
	    senderInfo: SenderInfo;
	    megSeq: number;
	    cts: number;
	    state: number;
	
	    static createFrom(source: any = {}) {
	        return new Msg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.msgId = source["msgId"];
	        this.convId = source["convId"];
	        this.sender = source["sender"];
	        this.receiver = source["receiver"];
	        this.chatType = source["chatType"];
	        this.msgType = source["msgType"];
	        this.content = this.convertValues(source["content"], MsgContent);
	        this.senderInfo = this.convertValues(source["senderInfo"], SenderInfo);
	        this.megSeq = source["megSeq"];
	        this.cts = source["cts"];
	        this.state = source["state"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class MsgSendData {
	    receiver?: string;
	    chatType?: number;
	    ttl?: number;
	    msgContent?: MsgContent;
	
	    static createFrom(source: any = {}) {
	        return new MsgSendData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.receiver = source["receiver"];
	        this.chatType = source["chatType"];
	        this.ttl = source["ttl"];
	        this.msgContent = this.convertValues(source["msgContent"], MsgContent);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

