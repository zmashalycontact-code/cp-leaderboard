import StatsDashboard from './StatsDashboard';
import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { 
  Trophy, Flame, AlertCircle, Smartphone, RotateCcw, 
  Target, Skull, Activity 
} from 'lucide-react';

const getRankColor = (tier) => {
  switch (tier) {
    case 'كحيان':         return 'text-gray-400 border-gray-600 bg-gray-400/10';
    case 'روش':           return 'text-cyan-400 border-cyan-600 bg-cyan-400/10';
    case 'باشا ستراكشر':  return 'text-teal-400 border-teal-600 bg-teal-400/10 shadow-[0_0_10px_rgba(45,212,191,0.3)]';
    case 'شكسبير':        return 'text-purple-400 border-purple-600 bg-purple-400/10 shadow-[0_0_10px_rgba(192,132,252,0.4)]';
    case 'تنين مجنح':     return 'text-red-500 border-red-600 bg-red-500/10 shadow-[0_0_15px_rgba(239,68,68,0.5)]';
    case 'The GOAT':      return 'text-fuchsia-400 border-fuchsia-500 bg-fuchsia-400/10 shadow-[0_0_20px_rgba(192,38,211,0.6)] animate-pulse';
    case 'CP MASTER':     return 'text-amber-400 border-amber-500 bg-amber-400/10 shadow-[0_0_25px_rgba(251,191,36,0.7)] font-black uppercase animate-pulse';
    default:              return 'text-slate-400 border-slate-700 bg-slate-800/50';
  }
};

const TableRow = ({ user, rank }) => {
  return (
    <motion.div 
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: rank * 0.05 }}
      className="grid grid-cols-[50px_minmax(150px,1fr)_120px_140px_100px_100px_100px_100px_100px] gap-4 py-4 px-6 items-center hover:bg-white/[0.03] transition-all border-b border-white/5 group relative overflow-hidden"
    >
      <div className="font-mono text-slate-500 font-bold">#{rank}</div>
      
      <div className="flex flex-col min-w-0 text-right">
        <span className="text-white font-bold truncate text-sm">{user.display_name}</span>
        <span className="text-[10px] text-slate-500 font-mono truncate">@{user.handle}</span>
        

        <div className="flex flex-wrap gap-2 mt-2">
          {user.peak_weekly_rating > 0 && (
            <span className="flex items-center gap-1 px-2 py-0.5 rounded-md text-[9px] font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20 uppercase tracking-tighter">
              <Target size={10} strokeWidth={3} /> Peak {user.peak_weekly_rating}
            </span>
          )}

          {user.struggle_count > 0 && (
            <span className="flex items-center gap-1 px-2 py-0.5 rounded-md text-[9px] font-bold bg-red-500/10 text-red-400 border border-red-500/20 uppercase tracking-tighter">
              <Skull size={10} strokeWidth={3} />
              <span>عنيد</span>

              <span className="bg-red-500/20 px-1 rounded ml-0.5">{user.struggle_count}</span>
            </span>
          )}
        </div>
      </div>

      <div className="flex justify-center">

        <span className={`px-2.5 py-1 rounded-md text-[9px] font-black tracking-widest uppercase border ${getRankColor(user.rank_tier)}`}>
          {user.rank_tier === 'CP MASTER' ? user.display_name : (user.rank_tier || 'UNRANKED')}
        </span>
      </div>

      <div className="text-center font-black text-white text-xl drop-shadow-[0_0_8px_rgba(255,255,255,0.3)]">
        {(user.season_points || 0).toFixed(1)}
      </div>

      <div className="text-center text-slate-400 font-mono font-bold">{(user.cf_points || 0).toFixed(1)}</div>
      

      <div className="text-center text-slate-400 font-mono font-bold">{(user.atcoder_points || 0).toFixed(1)}</div>
      
      <div className="text-center font-bold">
        {user.current_rating > 0 ? (
          <span className="text-orange-400 font-mono">{user.current_rating}</span>
        ) : (
          <span className="text-slate-700 text-[10px] italic font-bold uppercase">Unrated</span>
        )}
      </div>

      <div className="text-center text-slate-300 font-mono font-bold">{user.total_solved}</div>

      <div className="text-center text-emerald-400 font-mono font-black flex flex-col items-center justify-center gap-1">
        <div className="flex items-center gap-1">
          {user.activity_7d} <Flame size={12} fill="currentColor" className="opacity-80" />
        </div>
        {user.hidden_solved > 0 && (
          <motion.span 
            initial={{ scale: 0 }} animate={{ scale: 1 }}
            className="bg-indigo-500/20 text-indigo-400 text-[9px] font-bold px-2 py-0.5 rounded-full border border-indigo-500/30 whitespace-nowrap"
          >
            +{user.hidden_solved} شيتات
          </motion.span>
        )}
      </div>
    </motion.div>
  );
};

export default function App() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('table');

  useEffect(() => {

    const API_URL = import.meta.env.VITE_API_URL || 'https://zmashaly-idc-icpc.hf.space';
    
    fetch(`${API_URL}/api/leaderboard`)
      .then(res => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then(data => {

        const arr = Array.isArray(data)
          ? data
          : Object.values(data).map(i => (typeof i === 'string' ? JSON.parse(i) : i));
        
        setUsers(arr.filter(u => u && u.handle).sort((a, b) => (b.season_points || 0) - (a.season_points || 0)));
      })
      .catch(err => {
        console.error("Leaderboard fetch failed:", err);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="min-h-screen bg-[#020617] text-slate-200 font-sans p-4 md:p-8 selection:bg-indigo-500/30" dir="rtl">
      <div className="max-w-7xl mx-auto">
        
        <header className="mb-12 md:mb-20 text-center relative pt-10">
          <motion.div initial={{ y: -50, opacity: 0 }} animate={{ y: 0, opacity: 1 }}>
            <Trophy size={60} className="text-amber-400 mb-6 mx-auto animate-bounce drop-shadow-[0_0_20px_rgba(251,191,36,0.4)]" />
            <h1 className="text-4xl md:text-7xl font-black tracking-tighter text-white mb-4 italic uppercase">
              IDC CP <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400">Leaderboard</span>
            </h1>
            <p className="text-slate-500 font-bold tracking-widest uppercase text-[10px] md:text-xs text-center">
              Season 2026 • Delta University ICPC Community
            </p>
          </motion.div>
        </header>

        <div className="flex justify-center mb-10 gap-2">
          <button onClick={() => setActiveTab('table')} className={`px-8 py-3 rounded-full font-black text-xs uppercase tracking-[0.2em] transition-all ${activeTab === 'table' ? 'bg-indigo-600 text-white shadow-lg scale-105' : 'bg-white/5 text-slate-500 hover:bg-white/10'}`}>
            LEADERBOARD
          </button>
          <button onClick={() => setActiveTab('stats')} className={`px-8 py-3 rounded-full font-black text-xs uppercase tracking-[0.2em] transition-all ${activeTab === 'stats' ? 'bg-indigo-600 text-white shadow-lg scale-105' : 'bg-white/5 text-slate-500 hover:bg-white/10'}`}>
            STATS
          </button>
        </div>

        {loading ? (
          <div className="text-center py-20 animate-pulse text-xl font-black text-slate-500 uppercase tracking-widest">
            جاري مزامنة المصفوفة...
          </div>
        ) : (
          activeTab === 'table' ? (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              
              <div className="md:hidden portrait:flex landscape:hidden items-center justify-center gap-3 bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 p-3 rounded-2xl mb-6 shadow-lg animate-pulse">
                <Smartphone className="w-5 h-5 animate-bounce" />
                <span className="text-[10px] font-black italic text-center">لف الموبايل بالعرض أو استخدم اللابتوب لأفضل رؤية</span>
                <RotateCcw className="w-4 h-4" />
              </div>

              <div className="overflow-x-auto pb-6 custom-scrollbar">
                <div className="inline-block min-w-full align-middle px-2">
                  <div className="min-w-[1000px] bg-white/[0.02] border border-white/10 p-6 rounded-3xl shadow-2xl relative z-10 overflow-hidden">
                    <div className="grid grid-cols-[50px_minmax(150px,1fr)_120px_140px_100px_100px_100px_100px_100px] gap-4 py-4 px-6 text-[10px] font-black text-slate-600 uppercase tracking-widest border-b border-white/5 mb-6 text-center">
                      <div className="text-right">#</div>
                      <div className="text-right">المنافس</div>
                      <div>الرتبة</div>
                      <div>إجمالي النقاط</div>
                      <div>CF PTS</div>
                      <div>AC PTS</div>
                      <div className="text-orange-400/80">الريت</div>
                      <div>مسائل خام</div>
                      <div className="text-emerald-500">فورسز 7D</div>
                    </div>
                    {users.map((u, i) => <TableRow key={u.handle} user={u} rank={i + 1} />)}
                  </div>
                </div>
              </div>
            </motion.div>
          ) : (
            <StatsDashboard users={users} />
          )
        )}

        <footer className="mt-12 md:mt-20 pb-10 border-t border-white/5 pt-8 text-center text-slate-500 text-sm font-medium">
          Made with ❤️ by <span className="text-slate-300 font-bold">Ziad Mashaly</span>
        </footer>
      </div>
    </div>
  );
}
