/**
 * 将毫秒时间戳格式化为 IM 会话列表中友好的时间显示
 * @param {number} timestamp - 毫秒时间戳
 * @returns {string}
 */
export const prettyTime = (timestamp) => {
    if (!timestamp || typeof timestamp !== 'number') {
        return '--';
    }

    const msgDate = new Date(timestamp);
    const now = new Date();

    // 判断是否为同一天
    const isSameDay = (d1, d2) =>
        d1.getFullYear() === d2.getFullYear() &&
        d1.getMonth() === d2.getMonth() &&
        d1.getDate() === d2.getDate();

    // 1. 今天 → HH:mm
    if (isSameDay(msgDate, now)) {
        return padTime(msgDate.getHours()) + ':' + padTime(msgDate.getMinutes());
    }

    // 2. 昨天 → "昨天"
    const yesterday = new Date(now);
    yesterday.setDate(now.getDate() - 1);
    if (isSameDay(msgDate, yesterday)) {
        return '昨天';
    }

    // 3. 近7天内（且今年）→ 周几
    const sevenDaysAgo = new Date(now);
    sevenDaysAgo.setDate(now.getDate() - 7);
    if (
        msgDate >= sevenDaysAgo &&
        msgDate < yesterday &&
        msgDate.getFullYear() === now.getFullYear()
    ) {
        return weekday(msgDate.getDay());
    }

    // 4. 今年更早 → MM-dd
    if (msgDate.getFullYear() === now.getFullYear()) {
        return formatMMDD(msgDate);
    }

    // 5. 非今年 → yyyy-MM-dd
    return formatYYYYMMDD(msgDate);
};

// 工具函数：补零
const padTime = (n) => String(n).padStart(2, '0');

// 工具函数：格式化 MM-dd
const formatMMDD = (date) =>
    padTime(date.getMonth() + 1) + '-' + padTime(date.getDate());

// 工具函数：格式化 yyyy-MM-dd
const formatYYYYMMDD = (date) =>
    date.getFullYear() + '-' + formatMMDD(date);

// 工具函数：获取中文周几
const weekday = (day) => ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][day];