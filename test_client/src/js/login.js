var server = "http://localhost:4000"
$(document).ready(function() {
    $("#reg-tab").click(function(){
        $("#login-tab").removeClass("active-tab");
        $(this).addClass("active-tab");
        $("#login").hide();
        $("#register").show();
        $("#errors").text("");        
    })

    $("#login-tab").click(function(){
        $("#reg-tab").removeClass("active-tab");
        $(this).addClass("active-tab");
        $("#register").hide();
        $("#login").show(); 
        $("#errors").text("");
    })

    $("#login-sub").click(function(e) {
        e.preventDefault();
        var e = $("#log-email").val()
        var p = $("#log-pass").val()
        $.ajax({
            type: "POST",
            url: server+"/v1/sessions",
            data: JSON.stringify({"email": e,"password": p}),
            dataType: 'json',
            contentType: "application/json",
        }).done(function(data, status, xhr){
            var auth = xhr.getResponseHeader("Authorization")
            localStorage.setItem("auth", auth)
            window.location="./test.html"
        }).fail(function(xhr, status, error){
            $("#errors").text(xhr.responseText)
        })
    })

    $("#register-sub").click(function(e) {
        e.preventDefault();
        var e = $("#reg-email").val()
        var u = $("#reg-user").val()
        var fn = $("#reg-fn").val()
        var ln = $("#reg-ln").val()
        var p = $("#reg-pass").val()
        var pc = $("#reg-passconf").val()
        var payload = {
            "email": e,
            "username": u,
            "firstName": fn,
            "lastName": ln,
            "password": p,
            "passwordConf": pc
        }
        $.ajax({
            type: "POST",
            url: server+"/v1/users",
            data: JSON.stringify(payload),
            dataType: 'json',
            contentType: "application/json",
        }).done(function(data, status, xhr){
            var auth = xhr.getResponseHeader("Authorization")
            localStorage.setItem("auth", auth)
            window.location="./test.html"
        }).fail(function(xhr, status, error){
            $("#errors").text(xhr.responseText)
        })
    })
})

