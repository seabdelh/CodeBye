package com.herokuapp.codebye;

import android.app.Activity;
import android.content.Context;
import android.graphics.Color;
import android.text.Html;
import android.text.method.LinkMovementMethod;
import android.util.Log;
import android.view.Gravity;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseAdapter;
import android.widget.LinearLayout;
import android.widget.TextView;

import java.util.ArrayList;

public class ChatAdapter extends BaseAdapter{

    private static LayoutInflater inflater = null;
    ArrayList<ChatMessage> chatList;

    public ChatAdapter(Activity activity, ArrayList<ChatMessage> list) {
        chatList = list;
        inflater = (LayoutInflater) activity
                .getSystemService(Context.LAYOUT_INFLATER_SERVICE);

    }

    @Override
    public int getCount() {
        //Log.i("info","getCount");
        return chatList.size();
    }

    @Override
    public Object getItem(int i) {
       // Log.i("info","getItem");
        return i;
    }

    @Override
    public long getItemId(int i) {
       // Log.i("info","getItemId");
        return i;
    }

    @Override
    public View getView(int i, View CurrentView, ViewGroup viewGroup) {
        ChatMessage message = (ChatMessage) chatList.get(i);  //i is a position
        View view = CurrentView;
        if (view == null)
            view = inflater.inflate(R.layout.chatbubble, null);

        TextView msg = (TextView) view.findViewById(R.id.message_text);
        msg.setText(message.body);
        LinearLayout layout = (LinearLayout) view
                .findViewById(R.id.bubble_layout);
        LinearLayout parent_layout = (LinearLayout) view
                .findViewById(R.id.bubble_layout_parent);

        if (message.isMine) {
            layout.setBackgroundResource(R.drawable.bubble2);
            parent_layout.setGravity(Gravity.RIGHT);
        }
        else {
            layout.setBackgroundResource(R.drawable.bubble1);
            parent_layout.setGravity(Gravity.LEFT);
        }
        msg.setTextColor(Color.BLUE);
         if (message.body.startsWith("http")){
             Log.i("LINK","link sha8al");
             msg.setText(Html.fromHtml("<a href="+message.body+">"+ message.body));
             //Html.fromHtml(String, int) > 24

             msg.setMovementMethod(LinkMovementMethod.getInstance());
         }

        if(!message.isMine)
            msg.setTextColor(Color.BLACK);
        return view;

    }

    public void add(ChatMessage msg) {
        chatList.add(msg);
    }

}
