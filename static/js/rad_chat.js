// Wrap file in self calling function so we don't pollute global namespace
(function() {
'use strict;'

// A simple jQuery extension to check if an element exists
// https://stackoverflow.com/a/920322/2516916
$.fn.exists = function () {
    return this.length !== 0;
}

// Freaking IE11 caches GET API calls
// https://stackoverflow.com/a/4303862/2516916
$.ajaxSetup({ cache: false });

function insertMessage(data) {
   let m = $("#card_template").clone();

   // Unhide template copy
   m.removeAttr('hidden');

   // Set id to be message if
   m.attr("id", data.id);

   // Fill in message text
   m.find(".content").text(data.text);

    // Make upvote icon clickable
    m.find(".icon").click(function() {
       vote(data.id)
    });

   // Add to DOM
   $("#chat_messages").append(m);

   // Update upvote count next to icon
   updateUpvotes(data.id, data.upvotes);
}

function submitMessage() {
    let payload = {
        text: $("#chat_message").val()
    };

    $.post("/api/message", JSON.stringify(payload), function(data) {
        insertMessage(data);
        $("#chat_message").val(""); // clear chat box
    });
}

function updateUpvotes(id, upvotes) {
    let m = $("#"+id).find(".upvotes");

    m.text(upvotes);
    if (upvotes === 0)
        m.text(""); // Don't show a number if it is zero
}

function vote(id) {
    params = {
        id: id
    };

    // Find the icon element inside of the div
    let icon = $("#"+id).find("i");

    // If user already voted...
    if (icon.hasClass("fas")) {
        params.direction = "down"; // Unvote
        icon.toggleClass("fas far");
    }
    else {
        icon.toggleClass("far fas"); // Toggle arrow icon to be filled in
    }

    $.get("/api/vote", params, function(data) {
        updateUpvotes(id, data.upvotes);
    });
}

// Blindly populate the DOM with all messages
function populateMessages() {
    $.get("/api/messages", function(data) {
        $.each(data, function(i, value) {
            insertMessage(value);
        });
    });
}

let last_updated = new Date();
function updateMessages(last_updated) {
    let params = {
        updated_after: last_updated.toISOString()
    }

    $.get("/api/messages", params, function(data) {
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
    // Bind submit click with jQuery instead of `onclick=` because this file is wrapped in a function
    $("#chat_submit").click(submitMessage);

    populateMessages();
    setInterval(updateMessages, 5000, last_updated); // 5s
}

init()

})();