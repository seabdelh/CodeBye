package com.herokuapp.codebye;

import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;

import layout.Chat;

public class MainActivity extends AppCompatActivity implements Chat.OnFragmentInteractionListener {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

    }

    public void addMessage(ChatMessage chatMsg){

    }
}
