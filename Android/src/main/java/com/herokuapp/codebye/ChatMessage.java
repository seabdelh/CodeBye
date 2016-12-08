package com.herokuapp.codebye;

public class ChatMessage {

    public String body;
    public String Date, Time;
    public boolean isMine;

    public ChatMessage(String messageString,
                       boolean isMINE) {
        body = messageString;
        isMine = isMINE;
    }
}