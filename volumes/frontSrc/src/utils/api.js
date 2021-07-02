import request from '@/utils/request'

export async function getDirectorSelection() {
    const response = await request({
        url: '/stockitems',
        method: 'get'
    })
    
    response.map( (item , key) => {
        
        let closePrice = parseFloat(item.ClosePrice);
        let directorAvgPrice = parseFloat(item.DirectorSelection.AvgCost);
        item.ClosePrice = closePrice
        if(isNaN(directorAvgPrice)) {
            directorAvgPrice = 0
            item.DirectorSelection.AvgCost = directorAvgPrice
            item.DirectorSelection.CompareIsUp = false
            item.DirectorSelection.Percent = 0
        } else {
            item.DirectorSelection.AvgCost = directorAvgPrice
            item.DirectorSelection.CompareIsUp = closePrice > directorAvgPrice ? true : false
            if(isNaN(closePrice)) {
                item.ClosePrice = 0
                item.DirectorSelection.Percent = 0
            } else {
                item.DirectorSelection.Percent = parseFloat( (((closePrice - directorAvgPrice) / directorAvgPrice) * 100 ).toFixed(2) )
            }
        }
    })

    return response
}

export async function sendMmeberSpreadsheet(data){
    const response = await request({
        url: '/mss',
        method: 'post',
        data
    })
    
    response.map( (item , key) => {
        
        let closePrice = parseFloat(item.ClosePrice);
        item.ClosePrice = closePrice
        
        if(item.DirectorSelection !== null) {
            let directorAvgPrice = parseFloat(item.DirectorSelection.AvgCost);

            if(isNaN(directorAvgPrice)) {
                directorAvgPrice = 0
                item.DirectorSelection.AvgCost = directorAvgPrice
                item.DirectorSelection.CompareIsUp = false
                item.DirectorSelection.Percent = 0
            } else {
                item.DirectorSelection.AvgCost = directorAvgPrice
                item.DirectorSelection.CompareIsUp = closePrice > directorAvgPrice ? true : false
                if(isNaN(closePrice)) {
                    item.ClosePrice = 0
                    item.DirectorSelection.Percent = 0
                } else {
                    item.DirectorSelection.Percent = parseFloat( (((closePrice - directorAvgPrice) / directorAvgPrice) * 100 ).toFixed(2) )
                }
            }
        } else {
            item.DirectorSelection = {
                Percent : "-"
            }
        }

        
        if(item.MemberData !== null) {
            // 看一下平均成本是否正常
            item.MemberData.MemAvgCost = parseFloat(item.MemberData.MemAvgCost);
            // 如不正常
            if(isNaN(item.MemberData.MemAvgCost)) {
                item.MemberData.MemAvgCost = 0
                item.MemberData.Percent = 0
            } else { //正常計算獲利百分比
                item.MemberData.Percent = parseFloat( (((closePrice - item.MemberData.MemAvgCost) / item.MemberData.MemAvgCost) * 100 ).toFixed(2) )
                //計算獲利金額

                //目前價值                
                let nowValue = (closePrice * parseInt(item.MemberData.MemShares));
                let beforeValue = (item.MemberData.MemAvgCost * parseInt(item.MemberData.MemShares));
                let nowSaleFee = nowValue * 0.001425 //賣手續
                let beforeBuyFee = beforeValue * 0.001425 // 買手續
                // 計算最少手續費
                let nowValueDisc = nowSaleFee * parseFloat(item.MemberData.MemFeeDiscount);
                nowSaleFee = nowValueDisc > parseInt(item.MemberData.MemLowestFee) ? parseFloat(nowValueDisc.toFixed(2)) : parseInt(item.MemberData.MemLowestFee);
                
                let beforeValueDisc = beforeBuyFee * parseFloat(item.MemberData.MemFeeDiscount);
                beforeBuyFee = beforeValueDisc > parseInt(item.MemberData.MemLowestFee) ? parseFloat(beforeValueDisc.toFixed(2)) : parseInt(item.MemberData.MemLowestFee);                
                let saleTax = parseFloat( (nowValue * 0.003).toFixed(2) ) //證交稅
                // 這支股票收益
                let Profit = nowValue - beforeValue - (nowSaleFee + beforeBuyFee + saleTax);
                item.MemberData.Profit = parseFloat( Profit.toFixed(2) );
                // 這支股票投資額
                item.MemberData.CostValue = parseInt(beforeValue)
            }
            item.MemberData.CompareIsUp = closePrice > item.MemberData.MemAvgCost ? true : false
            
        } else {
            item.MemberData = {
                Percent : "-"
            }
        }
    })

    return response

    
}