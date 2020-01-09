package com.hiroebe.osushi

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle

import go.Seq

import com.hiroebe.osushi.mobile.EbitenView

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        Seq.setContext(applicationContext)
    }

    private fun getEbitenView(): EbitenView {
        return this.findViewById(R.id.ebitenview)
    }

    override fun onResume() {
        super.onResume()
        this.getEbitenView().resumeGame()
    }

    override fun onPause() {
        super.onPause()
        this.getEbitenView().suspendGame()
    }
}
