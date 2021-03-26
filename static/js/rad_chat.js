
'use strict;'

// https://stackoverflow.com/a/920322/2516916
$.fn.exists = function () {
    return this.length !== 0;
}

function insertMessage(data) {
   let m = $("#card_template").clone();
   m.removeAttr('hidden');
   m.attr("id", data.id);
   m.find(".content").text(data.text);
   m.find(".icon").attr("onclick", "vote(" + data.id + ")");

   $("#chat_messages").append(m);
   updateUpvotes(data.id, data.upvotes);
}

function submitMessage() {
    let payload = {
        text: $("#chat_message").val()
    };

    $.post("/api/message", JSON.stringify(payload), function(data) {
        console.log(data);
        insertMessage(data);
        $("#chat_message").val(""); // clear chat box
    });
}

function updateUpvotes(id, upvotes) {
    let m = $("#"+id).find(".upvotes");

    m.text(upvotes);
    if (upvotes === 0)
        m.text("");
}

function vote(id) {
    params = {
        id: id
    };

    let icon = $("#"+id).find("i");

    if (icon.hasClass("fas")) {
        params.direction = "down";
        icon.toggleClass("fas far");
    }
    else {
        icon.toggleClass("far fas");
    }

    $.get("/api/vote", params, function(data) {
        updateUpvotes(id, data.upvotes);
    });
}

function populateMessages() {
    $.get("/api/messages", function(data) {
        $.each(data, function(i, value) {
            insertMessage(value);
        });
    });
}

let last_updated = new Date();
function updateMessages() {
    let params = {
        updated_after: last_updated.toISOString()
    }

    $.get("/api/messages", params, function(data) {
        console.log(data)
        if (!data)
            return;

        $.each(data, function(i, value) {
            let update_time = new Date(value.last_updated);
            if (update_time > last_updated)
                last_updated = update_time;

            if ($("#"+value.id).exists()) {
                updateUpvotes(value.id, value.upvotes)
            }
            else {
                insertMessage(value);
            }

        });
    });
}

function init() {
    populateMessages();
    setInterval(updateMessages, 5000); // 5s
}

init()