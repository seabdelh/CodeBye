package layout;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.ViewGroup;
import android.content.Context;
import android.view.View;
import android.view.inputmethod.InputMethodManager;
import android.widget.AdapterView;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.ListView;

import com.android.volley.Request;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.VolleyLog;
import com.android.volley.toolbox.JsonObjectRequest;
import com.herokuapp.codebye.AppController;
import com.herokuapp.codebye.ChatAdapter;
import com.herokuapp.codebye.ChatMessage;
import com.herokuapp.codebye.R;

import org.json.JSONObject;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

public class Chat extends Fragment {
    private EditText msg_edittext;
    private ListView msgListView;
    private String uuid = null;
    private ImageButton sendButton;
    public static ArrayList<ChatMessage> chatlist;
    private ChatAdapter chatAdapter;
    private String baseUrl;
    private String msg;
    private String tag_json_obj = "json_obj_req";
    private OnFragmentInteractionListener mListener; //imp for tunneling messages from activity to fragment

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        View view = inflater.inflate(R.layout.fragment_chat, container, false);

        msg_edittext = (EditText) view.findViewById(R.id.text_chatMSG);
        msgListView = (ListView) view.findViewById(R.id.listView_chat);
        sendButton = (ImageButton) view
                .findViewById(R.id.btn_send);

        sendButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                sendButtonClicked();
            }
        });

        msgListView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position,
                                    long id) {
                hideSoftKeyboard(getActivity());

            }
        });

        // ----Set autoscroll of listview when a new message arrives----//
        msgListView.setTranscriptMode(ListView.TRANSCRIPT_MODE_ALWAYS_SCROLL);
        msgListView.setStackFromBottom(true);

        chatlist = new ArrayList<ChatMessage>();

        chatAdapter = new ChatAdapter(getActivity(), chatlist);
        msgListView.setAdapter(chatAdapter);


        baseUrl = "http://codebye.herokuapp.com/";
        String api = "welcome";
        msg = "DEFAULT MSG:PROBABLY ERROR";

        JsonObjectRequest jsonObjReq = new JsonObjectRequest(Request.Method.GET,
                baseUrl + api, null,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        try {
                            msg = response.getString("message");
                            uuid = response.getString("uuid");
                        } catch (Exception e) {

                        }
                        putMessage(msg, false);
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                VolleyLog.d("Api error:GET", "Error: " + error.getMessage());
            }
        });
        AppController.getInstance().addToRequestQueue(jsonObjReq, tag_json_obj);



        return view;
    }

    @Override
    public void onAttach(Context context) {
        super.onAttach(context);
        if (context instanceof OnFragmentInteractionListener) {
            mListener = (OnFragmentInteractionListener) context;
        } else {
            throw new RuntimeException(context.toString()
                    + " must implement OnFragmentInteractionListener");
        }

    }

    public void sendButtonClicked() {
        String toBeSent = msg_edittext.getText().toString();
        putMessage(toBeSent, true); // TODO: Get from textEdit the message string
        msg_edittext.setText("");
        postRequest(toBeSent);

    }

    private void putMessage(String msg, boolean isMine) {
        if (msg.startsWith("First name")) {
            String arr[] = msg.split(", ");
            msg="";
            for (String subStr : arr) {
                msg+=subStr+"\n";
            }
        }
            chatAdapter.add(new ChatMessage(msg, isMine));
            chatAdapter.notifyDataSetChanged();

    }

    private void postRequest(String toBeSent) {
        Log.i("infoPost", "postRequest");
        JSONObject params = new JSONObject();
        try {
            params.put("message", toBeSent);
        } catch (Exception e) {

        }
        JsonObjectRequest req = new JsonObjectRequest(
                Request.Method.POST,
                baseUrl + "chat",
                params,
                new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        try {
                            msg = response.getString("message");
                        } catch (Exception e) {
                        }
                        putMessage(msg, false);
                    }
                }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                Log.i("msgInfo","ERROR:POST" );

            }
        }
        ) {
            public Map<String, String> getHeaders() {
                Map<String, String> mHeaders = new HashMap<>();
                mHeaders.put("Authorization", uuid);
                return mHeaders;
            }
        };
        AppController.getInstance().addToRequestQueue(req, tag_json_obj);


    }
    private void hideSoftKeyboard(Activity activity) {
        InputMethodManager inputMethodManager =
                (InputMethodManager) activity.getSystemService(
                        Activity.INPUT_METHOD_SERVICE);
        inputMethodManager.hideSoftInputFromWindow(
                activity.getCurrentFocus().getWindowToken(), 0);
    }

    public interface OnFragmentInteractionListener {
        public void addMessage(ChatMessage chatMsg);
    }
}
