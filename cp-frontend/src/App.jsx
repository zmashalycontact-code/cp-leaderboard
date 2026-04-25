import StatsDashboard from './StatsDashboard';
import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Trophy, Skull, Flame, AlertCircle, Target, Activity, Smartphone, RotateCcw } from 'lucide-react';

const getRankColor = (tier) => {
  switch (tier) {
    case 'كحيان':         return 'text-gray-400 border-gray-600 bg-gray-400/10';
    case 'روش':           return 'text-cyan-400 border-cyan-600 bg-cyan-400/10';
    case 'باشا ستراكشر':  return 'text-teal-400 border-teal-600 bg-teal-400/10 shadow-[0_0_10px_rgba(45,212,191,0.3)]';
    case 'شكسبير':        return 'text-purple-400 border-purple-600 bg-purple-400/10 shadow-[0_0_10px_rgba(192,132,252,0.4)]';
    case 'تنين مجنح':     return 'text-red-500 border-red-600 bg-red-500/10 shadow-[0_0_15px_rgba(239,68,68,0.5)]';
    case 'The GOAT':      return 'text-fuchsia-400 border-fuchsia-500 bg-fuchsia-400/10 shadow-[0_0_15px_rgba(232,121,249,0.5)] animate-pulse';
    case 'CP MASTER':     return 'text-yellow-400 border-yellow-500 bg-yellow-400/10 shadow-[0_0_20px_rgba(250,204,21,0.6)] animate-pulse font-extrabold';
    default:              return 'text-gray-300 border-gray-500 bg-gray-500/10';
  }
};

const getRankNumberColor = (rank) => {
  switch (rank) {
    case 1:  return 'text-yellow-400';
    case 2:  return 'text-slate-300';
    case 3:  return 'text-orange-500';
    default: return 'text-slate-600';
  }
};

const TableRow = ({ user, rank }) => (
  <motion.div
    initial={{ opacity: 0, y: 10 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ delay: rank * 0.05 }}
    className={`grid grid-cols-[80px_minmax(200px,1fr)_120px_140px_100px_100px_100px_100px_100px] gap-4 items-center py-4 px-6 transition-all rounded-xl mb-3 border min-w-[1040px] ${
      user.is_cheater 
        ? 'bg-[linear-gradient(90deg,rgba(13,40,60,0.8),rgba(0,100,150,0.4))] border-l-4 border-l-[#00b4d8] border-y-cyan-900/30 border-r-cyan-900/30 opacity-60 grayscale-[40%] pointer-events-none'
        : 'bg-white/5 hover:bg-white/[0.08] border-white/5'
    }`}
  >
    <div className="text-xl font-mono flex items-center gap-1">
      <span className="text-xs opacity-50 text-slate-600">#</span>
      <span className={`font-black ${getRankNumberColor(rank)}`}>
        {user.is_cheater ? "🥶" : rank}
      </span>
    </div>

    <div className="text-right">
      <div className={`text-lg font-bold tracking-tight ${user.is_cheater ? 'text-cyan-600 line-through' : 'text-slate-100'}`}>
        {user.display_name}
      </div>
      
      <div className="flex flex-wrap items-center gap-2">
        <div className="text-xs text-slate-500 font-mono opacity-80">@{user.handle}</div>
        {user.is_cheater && (
          <span className="text-[9px] font-bold bg-cyan-900/50 text-cyan-300 border border-cyan-500/30 px-1.5 py-0.5 rounded uppercase">
            مُجمد (cheater)
          </span>
        )}
      </div>
      
      <div className="flex flex-wrap gap-2 mt-2">
        {user.peak_weekly_rating > 0 && !user.is_cheater && (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded-md text-[9px] font-bold bg-blue-500/10 text-blue-400 border border-blue-500/20 uppercase tracking-tighter">
            <Target size={10} strokeWidth={3} /> Peak {user.peak_weekly_rating}
          </span>
        )}

        {user.struggle_count > 0 && !user.is_cheater && (
          <span className="flex items-center gap-1 px-2 py-0.5 rounded-md text-[9px] font-bold bg-red-500/10 text-red-400 border border-red-500/20 uppercase tracking-tighter">
            <Skull size={10} strokeWidth={3} />
            <span>عنيد</span>
            <span className="bg-red-500/20 px-1 rounded ml-0.5">{user.struggle_count}</span>
          </span>
        )}
      </div>
    </div>

    <div className="text-center">
      <span className={`px-2.5 py-1 rounded-md text-[9px] font-black tracking-widest uppercase border ${user.is_cheater ? 'text-cyan-700 border-cyan-800 bg-cyan-900/10' : getRankColor(user.rank_tier)}`}>
        {user.is_cheater ? 'FROZEN' : (user.rank_tier === 'CP MASTER' ? user.display_name : (user.rank_tier || 'UNRANKED'))}
      </span>
    </div>

    <div className={`text-center text-2xl font-mono font-black ${user.is_cheater ? 'text-cyan-800 line-through' : 'text-white'}`}>{(user.season_points || 0).toFixed(1)}</div>
    <div className={`text-center font-mono font-bold ${user.is_cheater ? 'text-cyan-800 line-through' : 'text-slate-400'}`}>{(user.cf_points || 0).toFixed(1)}</div>
    <div className={`text-center font-mono font-bold ${user.is_cheater ? 'text-cyan-800 line-through' : 'text-slate-400'}`}>{(user.atcoder_points || 0).toFixed(1)}</div>
    
    <div className="text-center">
      {user.current_rating > 0
        ? <span className={`${user.is_cheater ? 'text-cyan-800 line-through' : 'text-orange-400'} font-mono font-bold`}>{user.current_rating}</span>
        : <span className="text-slate-700 text-[10px] font-bold italic uppercase">Unrated</span>
      }
    </div>

    <div className={`text-center font-mono font-bold ${user.is_cheater ? 'text-cyan-800 line-through' : 'text-slate-300'}`}>{user.total_solved}</div>
    
    <div className={`text-center font-mono font-black flex flex-col items-center justify-center gap-1 ${user.is_cheater ? 'text-cyan-800' : 'text-emerald-400'}`}>
      <div className="flex items-center gap-1">
        {user.activity_7d} <Flame size={12} fill="currentColor" className="opacity-80" />
      </div>
      {user.hidden_solved > 0 && !user.is_cheater && (
        <motion.span 
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          className="bg-indigo-500/20 text-indigo-400 text-[9px] font-bold px-2 py-0.5 rounded-full border border-indigo-500/30 whitespace-nowrap"
        >
          +{user.hidden_solved} شيتات شهرياً
        </motion.span>
      )}
    </div>
  </motion.div>
);
export default function App() {
  const [users, setUsers]       = useState([]);
  const [loading, setLoading]   = useState(true);
  const [error, setError]       = useState(null);
  const [activeTab, setActiveTab] = useState('leaderboard');

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
        
	setUsers(arr.filter(u => u && u.handle).sort((a, b) => {

		  if (a.is_cheater && !b.is_cheater) return 1;

		  if (!a.is_cheater && b.is_cheater) return -1;

		  return (b.season_points || 0) - (a.season_points || 0);
		}));
	      })
	      .catch(err => {
		console.error('Leaderboard fetch failed:', err);
		setError('تعذّر تحميل البيانات. حاول مرة أخرى.');
	      })
	      .finally(() => setLoading(false));
	  }, []);

  return (
    <div className="min-h-screen bg-black text-slate-200 p-4 md:p-8 font-sans" dir="rtl">
      <div className="max-w-7xl mx-auto">

        <header className="mb-12 border-b border-white/10 pb-8 flex flex-col md:flex-row justify-between items-center md:items-end gap-6 text-center md:text-right">
          <div className="flex flex-col md:flex-row items-center gap-6">
            <img
              src="/logo.jpg"
              alt="IDC ICPC Logo"
              className="w-24 h-24 rounded-3xl object-cover border-2 border-white/10 shadow-2xl"
              onError={(e) => { e.target.style.display = 'none'; }}
            />
            <div>
              <h1 className="text-5xl md:text-6xl font-black text-white mb-1 tracking-tighter italic uppercase">IDC CPers</h1>
              <p className="text-slate-500 font-bold text-base md:text-lg">ICPC Delta University Community</p>
            </div>
          </div>
          <div className="bg-emerald-500/10 text-emerald-400 px-4 py-2 rounded-full border border-emerald-500/20 text-xs font-bold animate-pulse h-fit">
            LIVE SYNC ACTIVE
          </div>
        </header>

        <div className="flex flex-col md:flex-row justify-center gap-4 md:gap-6 mb-12">
          <button
            onClick={() => setActiveTab('leaderboard')}
            className={`flex items-center justify-center gap-3 w-full md:w-auto px-10 py-4 rounded-2xl font-black text-xs tracking-[0.2em] uppercase transition-all duration-500 ${
              activeTab === 'leaderboard' 
              ? 'bg-emerald-500 text-black shadow-[0_0_30px_rgba(16,185,129,0.3)] scale-100 md:scale-105' 
              : 'bg-white/5 text-slate-500 hover:bg-white/10 hover:text-slate-300'
            }`}
          >
            <Trophy size={18} strokeWidth={2.5} />
            LEADERBOARD
          </button>
          
          <button
            onClick={() => setActiveTab('stats')}
            className={`flex items-center justify-center gap-3 w-full md:w-auto px-10 py-4 rounded-2xl font-black text-xs tracking-[0.2em] uppercase transition-all duration-500 ${
              activeTab === 'stats' 
              ? 'bg-purple-600 text-white shadow-[0_0_30px_rgba(147,51,234,0.3)] scale-100 md:scale-105' 
              : 'bg-white/5 text-slate-500 hover:bg-white/10 hover:text-slate-300'
            }`}
          >
            <Activity size={18} strokeWidth={2.5} />
            STATS
          </button>
        </div>

        {loading && (
          <div className="text-center py-20 animate-pulse text-xl md:text-2xl font-black">
            جاري مزامنة المصفوفة...
          </div>
        )}

        {error && !loading && (
          <div className="flex flex-col md:flex-row items-center justify-center gap-3 py-20 text-red-400 text-center">
            <AlertCircle size={24} />
            <span className="text-lg md:text-xl font-bold">{error}</span>
          </div>
        )}

        {!loading && !error && (
          activeTab === 'leaderboard' ? (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              

              <div className="md:hidden portrait:flex landscape:hidden items-center justify-center gap-3 bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 p-3 rounded-2xl mb-6 shadow-lg animate-pulse">
                <Smartphone className="w-5 h-5 animate-bounce" />
                <span className="text-[10px] font-black italic text-center">لف الموبايل بالعرض أو استخدم اللابتوب لأفضل رؤية</span>
                <RotateCcw className="w-4 h-4" />
              </div>

              <div className="overflow-x-auto pb-6 scrollbar-hide md:scrollbar-default">
                <div className="inline-block min-w-full align-middle px-2">
                   <div className="min-w-[1040px]">
                    <div className="grid grid-cols-[80px_minmax(200px,1fr)_120px_140px_100px_100px_100px_100px_100px] gap-4 py-4 px-6 text-[10px] font-black text-slate-600 uppercase tracking-widest border-b border-white/5 mb-6">
                      <div>#</div>
                      <div className="text-right">المنافس</div>
                      <div className="text-center">الرتبة</div>
                      <div className="text-center">إجمالي النقاط</div>
                      <div className="text-center">CF PTS</div>
                      <div className="text-center">AC PTS</div>
                      <div className="text-center text-orange-400/80">الريت</div>
                      <div className="text-center">مسائل خام</div>
                      <div className="text-center text-emerald-500">فورسز 7D</div>
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

	<footer className="mt-12 md:mt-20 pb-10 border-t border-white/5 pt-8 text-center">
		  <p className="text-slate-500 text-sm font-medium tracking-wide">
		    Made with <span className="text-red-500 animate-pulse inline-block">❤️</span> by 
		    <span className="text-slate-300 font-bold ml-1">Ziad Mashaly</span>
		  </p>
		  <p className="text-slate-600 text-[10px] mt-1 uppercase tracking-[0.2em]">
		    ICPC Delta University Community President
		  </p>
		</footer>

	      </div>
	    </div>
	  );
}
