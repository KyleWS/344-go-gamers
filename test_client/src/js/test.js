var auth = localStorage.getItem("auth")
var server = "http://localhost:4000"
const websocket = new WebSocket("ws://localhost:4000/v1/game/ws?auth="+localStorage.getItem("auth"));
websocket.addEventListener("message", function(event) { 
    eJSON = JSON.parse(event.data)
    console.log(eJSON)
    if (eJSON.action == "drawing-phase") {
        var d = JSON.stringify({"type":"drawing","data":{"objects":[{"type":"path","originX":"left","originY":"top","left":177.77499389648438,"top":60.75,"width":495.01,"height":232,"fill":null,"stroke":"black","strokeWidth":10,"strokeDashArray":null,"strokeLineCap":"round","strokeLineJoin":"round","strokeMiterLimit":10,"scaleX":1,"scaleY":1,"angle":0,"flipX":false,"flipY":false,"opacity":1,"shadow":null,"visible":true,"clipTo":null,"backgroundColor":"","fillRule":"nonzero","globalCompositeOperation":"source-over","transformMatrix":null,"skewX":0,"skewY":0,"pathOffset":{"x":430.27999389648437,"y":181.75},"path":[["M",215.76499389648438,197.75],["Q",215.77499389648438,197.75,217.27499389648438,197.75],["Q",218.77499389648438,197.75,222.77499389648438,197.75],["Q",226.77499389648438,197.75,242.27499389648438,197.75],["Q",257.7749938964844,197.75,270.2749938964844,197.75],["Q",282.7749938964844,197.75,288.2749938964844,197.75],["Q",293.7749938964844,197.75,295.7749938964844,199.75],["Q",297.7749938964844,201.75,301.7749938964844,204.75],["Q",305.7749938964844,207.75,308.7749938964844,215.25],["Q",311.7749938964844,222.75,311.7749938964844,230.25],["Q",311.7749938964844,237.75,311.7749938964844,247.75],["Q",311.7749938964844,257.75,306.7749938964844,265.75],["Q",301.7749938964844,273.75,298.2749938964844,276.75],["Q",294.7749938964844,279.75,280.2749938964844,288.75],["Q",265.7749938964844,297.75,248.27499389648438,297.75],["Q",230.77499389648438,297.75,220.27499389648438,291.75],["Q",209.77499389648438,285.75,203.77499389648438,277.75],["Q",197.77499389648438,269.75,190.27499389648438,245.75],["Q",182.77499389648438,221.75,182.77499389648438,207.75],["Q",182.77499389648438,193.75,184.27499389648438,190.75],["Q",185.77499389648438,187.75,191.77499389648438,180.75],["Q",197.77499389648438,173.75,223.77499389648438,155.75],["Q",249.77499389648438,137.75,283.7749938964844,121.75],["Q",317.7749938964844,105.75,401.7749938964844,85.75],["Q",485.7749938964844,65.75,547.7749938964844,65.75],["Q",609.7749938964844,65.75,630.2749938964844,69.75],["Q",650.7749938964844,73.75,656.2749938964844,75.75],["Q",661.7749938964844,77.75,666.7749938964844,79.75],["Q",671.7749938964844,81.75,673.7749938964844,82.25],["Q",675.7749938964844,82.75,676.7749938964844,82.75],["L",677.7849938964844,82.75]]}]}})

        setTimeout(() => {
            $.ajax({
                type: "POST",
                url: server +"/v1/game/submit",
                headers: {"Authorization": auth},
                data: d,
                dataType: 'json',
                contentType: "application/json",
            }).done(function(data, status, xhr){
                return
            }).fail(function(xhr, status, error){
                $("#errors").text(xhr.responseText)
            })
        },2000)
    }
    else if (eJSON.action == "description-phase" || eJSON.action == "first-phase" ) {
        setTimeout(() => {
            var d = JSON.stringify({"data": $("#data").val(), "type": "description"})
            $.ajax({
                type: "POST",
                url: server + "/v1/game/submit",
                headers: {"Authorization": auth},
                data: d,
                dataType: 'json',
                contentType: "application/json",
            }).done(function(data, status, xhr){
                return
            }).fail(function(xhr, status, error){
                $("#errors").text(xhr.responseText)
            })
        },2000)
    }
});
websocket.addEventListener("close", function(event) { 
    console.log(event.reason)
});
