// Generated from /home/phlimy/Projects/Contrib/stack/components/ledger/pkg/machine/script/NumScript.g4 by ANTLR 4.9.2
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class NumScriptParser extends Parser {
	static { RuntimeMetaData.checkVersion("4.9.2", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, NEWLINE=5, WHITESPACE=6, MULTILINE_COMMENT=7, 
		LINE_COMMENT=8, VARS=9, META=10, SET_TX_META=11, SET_ACCOUNT_META=12, 
		PRINT=13, FAIL=14, SEND=15, SOURCE=16, FROM=17, MAX=18, DESTINATION=19, 
		TO=20, ALLOCATE=21, OP_ADD=22, OP_SUB=23, LPAREN=24, RPAREN=25, LBRACK=26, 
		RBRACK=27, LBRACE=28, RBRACE=29, EQ=30, TY_ACCOUNT=31, TY_ASSET=32, TY_NUMBER=33, 
		TY_MONETARY=34, TY_PORTION=35, TY_STRING=36, STRING=37, PORTION=38, REMAINING=39, 
		KEPT=40, BALANCE=41, SAVE=42, NUMBER=43, PERCENT=44, VARIABLE_NAME=45, 
		ACCOUNT=46, ASSET=47;
	public static final int
		RULE_monetary = 0, RULE_monetaryAll = 1, RULE_literal = 2, RULE_variable = 3, 
		RULE_expression = 4, RULE_allotmentPortion = 5, RULE_destinationInOrder = 6, 
		RULE_destinationAllotment = 7, RULE_keptOrDestination = 8, RULE_destination = 9, 
		RULE_sourceAccountOverdraft = 10, RULE_sourceAccount = 11, RULE_sourceInOrder = 12, 
		RULE_sourceMaxed = 13, RULE_source = 14, RULE_sourceAllotment = 15, RULE_valueAwareSource = 16, 
		RULE_statement = 17, RULE_type_ = 18, RULE_origin = 19, RULE_varDecl = 20, 
		RULE_varListDecl = 21, RULE_script = 22;
	private static String[] makeRuleNames() {
		return new String[] {
			"monetary", "monetaryAll", "literal", "variable", "expression", "allotmentPortion", 
			"destinationInOrder", "destinationAllotment", "keptOrDestination", "destination", 
			"sourceAccountOverdraft", "sourceAccount", "sourceInOrder", "sourceMaxed", 
			"source", "sourceAllotment", "valueAwareSource", "statement", "type_", 
			"origin", "varDecl", "varListDecl", "script"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'*'", "'allowing overdraft up to'", "'allowing unbounded overdraft'", 
			"','", null, null, null, null, "'vars'", "'meta'", "'set_tx_meta'", "'set_account_meta'", 
			"'print'", "'fail'", "'send'", "'source'", "'from'", "'max'", "'destination'", 
			"'to'", "'allocate'", "'+'", "'-'", "'('", "')'", "'['", "']'", "'{'", 
			"'}'", "'='", "'account'", "'asset'", "'number'", "'monetary'", "'portion'", 
			"'string'", null, null, "'remaining'", "'kept'", "'balance'", "'save'", 
			null, "'%'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, null, null, null, null, "NEWLINE", "WHITESPACE", "MULTILINE_COMMENT", 
			"LINE_COMMENT", "VARS", "META", "SET_TX_META", "SET_ACCOUNT_META", "PRINT", 
			"FAIL", "SEND", "SOURCE", "FROM", "MAX", "DESTINATION", "TO", "ALLOCATE", 
			"OP_ADD", "OP_SUB", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "LBRACE", 
			"RBRACE", "EQ", "TY_ACCOUNT", "TY_ASSET", "TY_NUMBER", "TY_MONETARY", 
			"TY_PORTION", "TY_STRING", "STRING", "PORTION", "REMAINING", "KEPT", 
			"BALANCE", "SAVE", "NUMBER", "PERCENT", "VARIABLE_NAME", "ACCOUNT", "ASSET"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "NumScript.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public NumScriptParser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	public static class MonetaryContext extends ParserRuleContext {
		public ExpressionContext asset;
		public ExpressionContext amt;
		public TerminalNode LBRACK() { return getToken(NumScriptParser.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(NumScriptParser.RBRACK, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public MonetaryContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_monetary; }
	}

	public final MonetaryContext monetary() throws RecognitionException {
		MonetaryContext _localctx = new MonetaryContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_monetary);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(46);
			match(LBRACK);
			setState(47);
			((MonetaryContext)_localctx).asset = expression(0);
			setState(48);
			((MonetaryContext)_localctx).amt = expression(0);
			setState(49);
			match(RBRACK);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class MonetaryAllContext extends ParserRuleContext {
		public ExpressionContext asset;
		public TerminalNode LBRACK() { return getToken(NumScriptParser.LBRACK, 0); }
		public TerminalNode RBRACK() { return getToken(NumScriptParser.RBRACK, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public MonetaryAllContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_monetaryAll; }
	}

	public final MonetaryAllContext monetaryAll() throws RecognitionException {
		MonetaryAllContext _localctx = new MonetaryAllContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_monetaryAll);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(51);
			match(LBRACK);
			setState(52);
			((MonetaryAllContext)_localctx).asset = expression(0);
			setState(53);
			match(T__0);
			setState(54);
			match(RBRACK);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class LiteralContext extends ParserRuleContext {
		public LiteralContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_literal; }
	 
		public LiteralContext() { }
		public void copyFrom(LiteralContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class LitPortionContext extends LiteralContext {
		public TerminalNode PORTION() { return getToken(NumScriptParser.PORTION, 0); }
		public LitPortionContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	public static class LitStringContext extends LiteralContext {
		public TerminalNode STRING() { return getToken(NumScriptParser.STRING, 0); }
		public LitStringContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	public static class LitAccountContext extends LiteralContext {
		public TerminalNode ACCOUNT() { return getToken(NumScriptParser.ACCOUNT, 0); }
		public LitAccountContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	public static class LitAssetContext extends LiteralContext {
		public TerminalNode ASSET() { return getToken(NumScriptParser.ASSET, 0); }
		public LitAssetContext(LiteralContext ctx) { copyFrom(ctx); }
	}
	public static class LitNumberContext extends LiteralContext {
		public TerminalNode NUMBER() { return getToken(NumScriptParser.NUMBER, 0); }
		public LitNumberContext(LiteralContext ctx) { copyFrom(ctx); }
	}

	public final LiteralContext literal() throws RecognitionException {
		LiteralContext _localctx = new LiteralContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_literal);
		try {
			setState(61);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case ACCOUNT:
				_localctx = new LitAccountContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(56);
				match(ACCOUNT);
				}
				break;
			case ASSET:
				_localctx = new LitAssetContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(57);
				match(ASSET);
				}
				break;
			case NUMBER:
				_localctx = new LitNumberContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(58);
				match(NUMBER);
				}
				break;
			case STRING:
				_localctx = new LitStringContext(_localctx);
				enterOuterAlt(_localctx, 4);
				{
				setState(59);
				match(STRING);
				}
				break;
			case PORTION:
				_localctx = new LitPortionContext(_localctx);
				enterOuterAlt(_localctx, 5);
				{
				setState(60);
				match(PORTION);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class VariableContext extends ParserRuleContext {
		public TerminalNode VARIABLE_NAME() { return getToken(NumScriptParser.VARIABLE_NAME, 0); }
		public VariableContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_variable; }
	}

	public final VariableContext variable() throws RecognitionException {
		VariableContext _localctx = new VariableContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_variable);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(63);
			match(VARIABLE_NAME);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ExpressionContext extends ParserRuleContext {
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
	 
		public ExpressionContext() { }
		public void copyFrom(ExpressionContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class ExprAddSubContext extends ExpressionContext {
		public ExpressionContext lhs;
		public Token op;
		public ExpressionContext rhs;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode OP_ADD() { return getToken(NumScriptParser.OP_ADD, 0); }
		public TerminalNode OP_SUB() { return getToken(NumScriptParser.OP_SUB, 0); }
		public ExprAddSubContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ExprLiteralContext extends ExpressionContext {
		public LiteralContext lit;
		public LiteralContext literal() {
			return getRuleContext(LiteralContext.class,0);
		}
		public ExprLiteralContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ExprVariableContext extends ExpressionContext {
		public VariableContext var_;
		public VariableContext variable() {
			return getRuleContext(VariableContext.class,0);
		}
		public ExprVariableContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	public static class ExprMonetaryNewContext extends ExpressionContext {
		public MonetaryContext mon;
		public MonetaryContext monetary() {
			return getRuleContext(MonetaryContext.class,0);
		}
		public ExprMonetaryNewContext(ExpressionContext ctx) { copyFrom(ctx); }
	}

	public final ExpressionContext expression() throws RecognitionException {
		return expression(0);
	}

	private ExpressionContext expression(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		ExpressionContext _localctx = new ExpressionContext(_ctx, _parentState);
		ExpressionContext _prevctx = _localctx;
		int _startState = 8;
		enterRecursionRule(_localctx, 8, RULE_expression, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(69);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case STRING:
			case PORTION:
			case NUMBER:
			case ACCOUNT:
			case ASSET:
				{
				_localctx = new ExprLiteralContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(66);
				((ExprLiteralContext)_localctx).lit = literal();
				}
				break;
			case VARIABLE_NAME:
				{
				_localctx = new ExprVariableContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(67);
				((ExprVariableContext)_localctx).var_ = variable();
				}
				break;
			case LBRACK:
				{
				_localctx = new ExprMonetaryNewContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(68);
				((ExprMonetaryNewContext)_localctx).mon = monetary();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			_ctx.stop = _input.LT(-1);
			setState(76);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,2,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					{
					_localctx = new ExprAddSubContext(new ExpressionContext(_parentctx, _parentState));
					((ExprAddSubContext)_localctx).lhs = _prevctx;
					pushNewRecursionContext(_localctx, _startState, RULE_expression);
					setState(71);
					if (!(precpred(_ctx, 4))) throw new FailedPredicateException(this, "precpred(_ctx, 4)");
					setState(72);
					((ExprAddSubContext)_localctx).op = _input.LT(1);
					_la = _input.LA(1);
					if ( !(_la==OP_ADD || _la==OP_SUB) ) {
						((ExprAddSubContext)_localctx).op = (Token)_errHandler.recoverInline(this);
					}
					else {
						if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
						_errHandler.reportMatch(this);
						consume();
					}
					setState(73);
					((ExprAddSubContext)_localctx).rhs = expression(5);
					}
					} 
				}
				setState(78);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,2,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	public static class AllotmentPortionContext extends ParserRuleContext {
		public AllotmentPortionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_allotmentPortion; }
	 
		public AllotmentPortionContext() { }
		public void copyFrom(AllotmentPortionContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class AllotmentPortionRemainingContext extends AllotmentPortionContext {
		public TerminalNode REMAINING() { return getToken(NumScriptParser.REMAINING, 0); }
		public AllotmentPortionRemainingContext(AllotmentPortionContext ctx) { copyFrom(ctx); }
	}
	public static class AllotmentPortionVarContext extends AllotmentPortionContext {
		public VariableContext por;
		public VariableContext variable() {
			return getRuleContext(VariableContext.class,0);
		}
		public AllotmentPortionVarContext(AllotmentPortionContext ctx) { copyFrom(ctx); }
	}
	public static class AllotmentPortionConstContext extends AllotmentPortionContext {
		public TerminalNode PORTION() { return getToken(NumScriptParser.PORTION, 0); }
		public AllotmentPortionConstContext(AllotmentPortionContext ctx) { copyFrom(ctx); }
	}

	public final AllotmentPortionContext allotmentPortion() throws RecognitionException {
		AllotmentPortionContext _localctx = new AllotmentPortionContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_allotmentPortion);
		try {
			setState(82);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case PORTION:
				_localctx = new AllotmentPortionConstContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(79);
				match(PORTION);
				}
				break;
			case VARIABLE_NAME:
				_localctx = new AllotmentPortionVarContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(80);
				((AllotmentPortionVarContext)_localctx).por = variable();
				}
				break;
			case REMAINING:
				_localctx = new AllotmentPortionRemainingContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(81);
				match(REMAINING);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class DestinationInOrderContext extends ParserRuleContext {
		public ExpressionContext expression;
		public List<ExpressionContext> amounts = new ArrayList<ExpressionContext>();
		public KeptOrDestinationContext keptOrDestination;
		public List<KeptOrDestinationContext> dests = new ArrayList<KeptOrDestinationContext>();
		public KeptOrDestinationContext remainingDest;
		public TerminalNode LBRACE() { return getToken(NumScriptParser.LBRACE, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode REMAINING() { return getToken(NumScriptParser.REMAINING, 0); }
		public TerminalNode RBRACE() { return getToken(NumScriptParser.RBRACE, 0); }
		public List<KeptOrDestinationContext> keptOrDestination() {
			return getRuleContexts(KeptOrDestinationContext.class);
		}
		public KeptOrDestinationContext keptOrDestination(int i) {
			return getRuleContext(KeptOrDestinationContext.class,i);
		}
		public List<TerminalNode> MAX() { return getTokens(NumScriptParser.MAX); }
		public TerminalNode MAX(int i) {
			return getToken(NumScriptParser.MAX, i);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public DestinationInOrderContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_destinationInOrder; }
	}

	public final DestinationInOrderContext destinationInOrder() throws RecognitionException {
		DestinationInOrderContext _localctx = new DestinationInOrderContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_destinationInOrder);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(84);
			match(LBRACE);
			setState(85);
			match(NEWLINE);
			setState(91); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(86);
				match(MAX);
				setState(87);
				((DestinationInOrderContext)_localctx).expression = expression(0);
				((DestinationInOrderContext)_localctx).amounts.add(((DestinationInOrderContext)_localctx).expression);
				setState(88);
				((DestinationInOrderContext)_localctx).keptOrDestination = keptOrDestination();
				((DestinationInOrderContext)_localctx).dests.add(((DestinationInOrderContext)_localctx).keptOrDestination);
				setState(89);
				match(NEWLINE);
				}
				}
				setState(93); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( _la==MAX );
			setState(95);
			match(REMAINING);
			setState(96);
			((DestinationInOrderContext)_localctx).remainingDest = keptOrDestination();
			setState(97);
			match(NEWLINE);
			setState(98);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class DestinationAllotmentContext extends ParserRuleContext {
		public AllotmentPortionContext allotmentPortion;
		public List<AllotmentPortionContext> portions = new ArrayList<AllotmentPortionContext>();
		public KeptOrDestinationContext keptOrDestination;
		public List<KeptOrDestinationContext> dests = new ArrayList<KeptOrDestinationContext>();
		public TerminalNode LBRACE() { return getToken(NumScriptParser.LBRACE, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RBRACE() { return getToken(NumScriptParser.RBRACE, 0); }
		public List<AllotmentPortionContext> allotmentPortion() {
			return getRuleContexts(AllotmentPortionContext.class);
		}
		public AllotmentPortionContext allotmentPortion(int i) {
			return getRuleContext(AllotmentPortionContext.class,i);
		}
		public List<KeptOrDestinationContext> keptOrDestination() {
			return getRuleContexts(KeptOrDestinationContext.class);
		}
		public KeptOrDestinationContext keptOrDestination(int i) {
			return getRuleContext(KeptOrDestinationContext.class,i);
		}
		public DestinationAllotmentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_destinationAllotment; }
	}

	public final DestinationAllotmentContext destinationAllotment() throws RecognitionException {
		DestinationAllotmentContext _localctx = new DestinationAllotmentContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_destinationAllotment);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(100);
			match(LBRACE);
			setState(101);
			match(NEWLINE);
			setState(106); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(102);
				((DestinationAllotmentContext)_localctx).allotmentPortion = allotmentPortion();
				((DestinationAllotmentContext)_localctx).portions.add(((DestinationAllotmentContext)_localctx).allotmentPortion);
				setState(103);
				((DestinationAllotmentContext)_localctx).keptOrDestination = keptOrDestination();
				((DestinationAllotmentContext)_localctx).dests.add(((DestinationAllotmentContext)_localctx).keptOrDestination);
				setState(104);
				match(NEWLINE);
				}
				}
				setState(108); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << PORTION) | (1L << REMAINING) | (1L << VARIABLE_NAME))) != 0) );
			setState(110);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class KeptOrDestinationContext extends ParserRuleContext {
		public KeptOrDestinationContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_keptOrDestination; }
	 
		public KeptOrDestinationContext() { }
		public void copyFrom(KeptOrDestinationContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class IsKeptContext extends KeptOrDestinationContext {
		public TerminalNode KEPT() { return getToken(NumScriptParser.KEPT, 0); }
		public IsKeptContext(KeptOrDestinationContext ctx) { copyFrom(ctx); }
	}
	public static class IsDestinationContext extends KeptOrDestinationContext {
		public TerminalNode TO() { return getToken(NumScriptParser.TO, 0); }
		public DestinationContext destination() {
			return getRuleContext(DestinationContext.class,0);
		}
		public IsDestinationContext(KeptOrDestinationContext ctx) { copyFrom(ctx); }
	}

	public final KeptOrDestinationContext keptOrDestination() throws RecognitionException {
		KeptOrDestinationContext _localctx = new KeptOrDestinationContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_keptOrDestination);
		try {
			setState(115);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case TO:
				_localctx = new IsDestinationContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(112);
				match(TO);
				setState(113);
				destination();
				}
				break;
			case KEPT:
				_localctx = new IsKeptContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(114);
				match(KEPT);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class DestinationContext extends ParserRuleContext {
		public DestinationContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_destination; }
	 
		public DestinationContext() { }
		public void copyFrom(DestinationContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class DestAccountContext extends DestinationContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public DestAccountContext(DestinationContext ctx) { copyFrom(ctx); }
	}
	public static class DestAllotmentContext extends DestinationContext {
		public DestinationAllotmentContext destinationAllotment() {
			return getRuleContext(DestinationAllotmentContext.class,0);
		}
		public DestAllotmentContext(DestinationContext ctx) { copyFrom(ctx); }
	}
	public static class DestInOrderContext extends DestinationContext {
		public DestinationInOrderContext destinationInOrder() {
			return getRuleContext(DestinationInOrderContext.class,0);
		}
		public DestInOrderContext(DestinationContext ctx) { copyFrom(ctx); }
	}

	public final DestinationContext destination() throws RecognitionException {
		DestinationContext _localctx = new DestinationContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_destination);
		try {
			setState(120);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,7,_ctx) ) {
			case 1:
				_localctx = new DestAccountContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(117);
				expression(0);
				}
				break;
			case 2:
				_localctx = new DestInOrderContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(118);
				destinationInOrder();
				}
				break;
			case 3:
				_localctx = new DestAllotmentContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(119);
				destinationAllotment();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceAccountOverdraftContext extends ParserRuleContext {
		public SourceAccountOverdraftContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sourceAccountOverdraft; }
	 
		public SourceAccountOverdraftContext() { }
		public void copyFrom(SourceAccountOverdraftContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class SrcAccountOverdraftSpecificContext extends SourceAccountOverdraftContext {
		public ExpressionContext specific;
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public SrcAccountOverdraftSpecificContext(SourceAccountOverdraftContext ctx) { copyFrom(ctx); }
	}
	public static class SrcAccountOverdraftUnboundedContext extends SourceAccountOverdraftContext {
		public SrcAccountOverdraftUnboundedContext(SourceAccountOverdraftContext ctx) { copyFrom(ctx); }
	}

	public final SourceAccountOverdraftContext sourceAccountOverdraft() throws RecognitionException {
		SourceAccountOverdraftContext _localctx = new SourceAccountOverdraftContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_sourceAccountOverdraft);
		try {
			setState(125);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case T__1:
				_localctx = new SrcAccountOverdraftSpecificContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(122);
				match(T__1);
				setState(123);
				((SrcAccountOverdraftSpecificContext)_localctx).specific = expression(0);
				}
				break;
			case T__2:
				_localctx = new SrcAccountOverdraftUnboundedContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(124);
				match(T__2);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceAccountContext extends ParserRuleContext {
		public ExpressionContext account;
		public SourceAccountOverdraftContext overdraft;
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public SourceAccountOverdraftContext sourceAccountOverdraft() {
			return getRuleContext(SourceAccountOverdraftContext.class,0);
		}
		public SourceAccountContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sourceAccount; }
	}

	public final SourceAccountContext sourceAccount() throws RecognitionException {
		SourceAccountContext _localctx = new SourceAccountContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_sourceAccount);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(127);
			((SourceAccountContext)_localctx).account = expression(0);
			setState(129);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==T__1 || _la==T__2) {
				{
				setState(128);
				((SourceAccountContext)_localctx).overdraft = sourceAccountOverdraft();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceInOrderContext extends ParserRuleContext {
		public SourceContext source;
		public List<SourceContext> sources = new ArrayList<SourceContext>();
		public TerminalNode LBRACE() { return getToken(NumScriptParser.LBRACE, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RBRACE() { return getToken(NumScriptParser.RBRACE, 0); }
		public List<SourceContext> source() {
			return getRuleContexts(SourceContext.class);
		}
		public SourceContext source(int i) {
			return getRuleContext(SourceContext.class,i);
		}
		public SourceInOrderContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sourceInOrder; }
	}

	public final SourceInOrderContext sourceInOrder() throws RecognitionException {
		SourceInOrderContext _localctx = new SourceInOrderContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_sourceInOrder);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(131);
			match(LBRACE);
			setState(132);
			match(NEWLINE);
			setState(136); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(133);
				((SourceInOrderContext)_localctx).source = source();
				((SourceInOrderContext)_localctx).sources.add(((SourceInOrderContext)_localctx).source);
				setState(134);
				match(NEWLINE);
				}
				}
				setState(138); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << MAX) | (1L << LBRACK) | (1L << LBRACE) | (1L << STRING) | (1L << PORTION) | (1L << NUMBER) | (1L << VARIABLE_NAME) | (1L << ACCOUNT) | (1L << ASSET))) != 0) );
			setState(140);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceMaxedContext extends ParserRuleContext {
		public ExpressionContext max;
		public SourceContext src;
		public TerminalNode MAX() { return getToken(NumScriptParser.MAX, 0); }
		public TerminalNode FROM() { return getToken(NumScriptParser.FROM, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public SourceContext source() {
			return getRuleContext(SourceContext.class,0);
		}
		public SourceMaxedContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sourceMaxed; }
	}

	public final SourceMaxedContext sourceMaxed() throws RecognitionException {
		SourceMaxedContext _localctx = new SourceMaxedContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_sourceMaxed);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(142);
			match(MAX);
			setState(143);
			((SourceMaxedContext)_localctx).max = expression(0);
			setState(144);
			match(FROM);
			setState(145);
			((SourceMaxedContext)_localctx).src = source();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceContext extends ParserRuleContext {
		public SourceContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_source; }
	 
		public SourceContext() { }
		public void copyFrom(SourceContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class SrcAccountContext extends SourceContext {
		public SourceAccountContext sourceAccount() {
			return getRuleContext(SourceAccountContext.class,0);
		}
		public SrcAccountContext(SourceContext ctx) { copyFrom(ctx); }
	}
	public static class SrcMaxedContext extends SourceContext {
		public SourceMaxedContext sourceMaxed() {
			return getRuleContext(SourceMaxedContext.class,0);
		}
		public SrcMaxedContext(SourceContext ctx) { copyFrom(ctx); }
	}
	public static class SrcInOrderContext extends SourceContext {
		public SourceInOrderContext sourceInOrder() {
			return getRuleContext(SourceInOrderContext.class,0);
		}
		public SrcInOrderContext(SourceContext ctx) { copyFrom(ctx); }
	}

	public final SourceContext source() throws RecognitionException {
		SourceContext _localctx = new SourceContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_source);
		try {
			setState(150);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case LBRACK:
			case STRING:
			case PORTION:
			case NUMBER:
			case VARIABLE_NAME:
			case ACCOUNT:
			case ASSET:
				_localctx = new SrcAccountContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(147);
				sourceAccount();
				}
				break;
			case MAX:
				_localctx = new SrcMaxedContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(148);
				sourceMaxed();
				}
				break;
			case LBRACE:
				_localctx = new SrcInOrderContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(149);
				sourceInOrder();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class SourceAllotmentContext extends ParserRuleContext {
		public AllotmentPortionContext allotmentPortion;
		public List<AllotmentPortionContext> portions = new ArrayList<AllotmentPortionContext>();
		public SourceContext source;
		public List<SourceContext> sources = new ArrayList<SourceContext>();
		public TerminalNode LBRACE() { return getToken(NumScriptParser.LBRACE, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RBRACE() { return getToken(NumScriptParser.RBRACE, 0); }
		public List<TerminalNode> FROM() { return getTokens(NumScriptParser.FROM); }
		public TerminalNode FROM(int i) {
			return getToken(NumScriptParser.FROM, i);
		}
		public List<AllotmentPortionContext> allotmentPortion() {
			return getRuleContexts(AllotmentPortionContext.class);
		}
		public AllotmentPortionContext allotmentPortion(int i) {
			return getRuleContext(AllotmentPortionContext.class,i);
		}
		public List<SourceContext> source() {
			return getRuleContexts(SourceContext.class);
		}
		public SourceContext source(int i) {
			return getRuleContext(SourceContext.class,i);
		}
		public SourceAllotmentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_sourceAllotment; }
	}

	public final SourceAllotmentContext sourceAllotment() throws RecognitionException {
		SourceAllotmentContext _localctx = new SourceAllotmentContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_sourceAllotment);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(152);
			match(LBRACE);
			setState(153);
			match(NEWLINE);
			setState(159); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(154);
				((SourceAllotmentContext)_localctx).allotmentPortion = allotmentPortion();
				((SourceAllotmentContext)_localctx).portions.add(((SourceAllotmentContext)_localctx).allotmentPortion);
				setState(155);
				match(FROM);
				setState(156);
				((SourceAllotmentContext)_localctx).source = source();
				((SourceAllotmentContext)_localctx).sources.add(((SourceAllotmentContext)_localctx).source);
				setState(157);
				match(NEWLINE);
				}
				}
				setState(161); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << PORTION) | (1L << REMAINING) | (1L << VARIABLE_NAME))) != 0) );
			setState(163);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ValueAwareSourceContext extends ParserRuleContext {
		public ValueAwareSourceContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_valueAwareSource; }
	 
		public ValueAwareSourceContext() { }
		public void copyFrom(ValueAwareSourceContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class SrcContext extends ValueAwareSourceContext {
		public SourceContext source() {
			return getRuleContext(SourceContext.class,0);
		}
		public SrcContext(ValueAwareSourceContext ctx) { copyFrom(ctx); }
	}
	public static class SrcAllotmentContext extends ValueAwareSourceContext {
		public SourceAllotmentContext sourceAllotment() {
			return getRuleContext(SourceAllotmentContext.class,0);
		}
		public SrcAllotmentContext(ValueAwareSourceContext ctx) { copyFrom(ctx); }
	}

	public final ValueAwareSourceContext valueAwareSource() throws RecognitionException {
		ValueAwareSourceContext _localctx = new ValueAwareSourceContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_valueAwareSource);
		try {
			setState(167);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,13,_ctx) ) {
			case 1:
				_localctx = new SrcContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(165);
				source();
				}
				break;
			case 2:
				_localctx = new SrcAllotmentContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(166);
				sourceAllotment();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class StatementContext extends ParserRuleContext {
		public StatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_statement; }
	 
		public StatementContext() { }
		public void copyFrom(StatementContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class PrintContext extends StatementContext {
		public ExpressionContext expr;
		public TerminalNode PRINT() { return getToken(NumScriptParser.PRINT, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public PrintContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class SendAllContext extends StatementContext {
		public MonetaryAllContext monAll;
		public SourceContext src;
		public DestinationContext dest;
		public TerminalNode SEND() { return getToken(NumScriptParser.SEND, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public MonetaryAllContext monetaryAll() {
			return getRuleContext(MonetaryAllContext.class,0);
		}
		public TerminalNode SOURCE() { return getToken(NumScriptParser.SOURCE, 0); }
		public List<TerminalNode> EQ() { return getTokens(NumScriptParser.EQ); }
		public TerminalNode EQ(int i) {
			return getToken(NumScriptParser.EQ, i);
		}
		public TerminalNode DESTINATION() { return getToken(NumScriptParser.DESTINATION, 0); }
		public SourceContext source() {
			return getRuleContext(SourceContext.class,0);
		}
		public DestinationContext destination() {
			return getRuleContext(DestinationContext.class,0);
		}
		public SendAllContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class SaveFromAccountContext extends StatementContext {
		public ExpressionContext mon;
		public MonetaryAllContext monAll;
		public ExpressionContext acc;
		public TerminalNode SAVE() { return getToken(NumScriptParser.SAVE, 0); }
		public TerminalNode FROM() { return getToken(NumScriptParser.FROM, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public MonetaryAllContext monetaryAll() {
			return getRuleContext(MonetaryAllContext.class,0);
		}
		public SaveFromAccountContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class SetTxMetaContext extends StatementContext {
		public Token key;
		public ExpressionContext value;
		public TerminalNode SET_TX_META() { return getToken(NumScriptParser.SET_TX_META, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public TerminalNode STRING() { return getToken(NumScriptParser.STRING, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public SetTxMetaContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class SetAccountMetaContext extends StatementContext {
		public ExpressionContext acc;
		public Token key;
		public ExpressionContext value;
		public TerminalNode SET_ACCOUNT_META() { return getToken(NumScriptParser.SET_ACCOUNT_META, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode STRING() { return getToken(NumScriptParser.STRING, 0); }
		public SetAccountMetaContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class FailContext extends StatementContext {
		public TerminalNode FAIL() { return getToken(NumScriptParser.FAIL, 0); }
		public FailContext(StatementContext ctx) { copyFrom(ctx); }
	}
	public static class SendContext extends StatementContext {
		public ExpressionContext mon;
		public ValueAwareSourceContext src;
		public DestinationContext dest;
		public TerminalNode SEND() { return getToken(NumScriptParser.SEND, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SOURCE() { return getToken(NumScriptParser.SOURCE, 0); }
		public List<TerminalNode> EQ() { return getTokens(NumScriptParser.EQ); }
		public TerminalNode EQ(int i) {
			return getToken(NumScriptParser.EQ, i);
		}
		public TerminalNode DESTINATION() { return getToken(NumScriptParser.DESTINATION, 0); }
		public ValueAwareSourceContext valueAwareSource() {
			return getRuleContext(ValueAwareSourceContext.class,0);
		}
		public DestinationContext destination() {
			return getRuleContext(DestinationContext.class,0);
		}
		public SendContext(StatementContext ctx) { copyFrom(ctx); }
	}

	public final StatementContext statement() throws RecognitionException {
		StatementContext _localctx = new StatementContext(_ctx, getState());
		enterRule(_localctx, 34, RULE_statement);
		try {
			setState(246);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,17,_ctx) ) {
			case 1:
				_localctx = new PrintContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(169);
				match(PRINT);
				setState(170);
				((PrintContext)_localctx).expr = expression(0);
				}
				break;
			case 2:
				_localctx = new SaveFromAccountContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(171);
				match(SAVE);
				setState(174);
				_errHandler.sync(this);
				switch ( getInterpreter().adaptivePredict(_input,14,_ctx) ) {
				case 1:
					{
					setState(172);
					((SaveFromAccountContext)_localctx).mon = expression(0);
					}
					break;
				case 2:
					{
					setState(173);
					((SaveFromAccountContext)_localctx).monAll = monetaryAll();
					}
					break;
				}
				setState(176);
				match(FROM);
				setState(177);
				((SaveFromAccountContext)_localctx).acc = expression(0);
				}
				break;
			case 3:
				_localctx = new SetTxMetaContext(_localctx);
				enterOuterAlt(_localctx, 3);
				{
				setState(179);
				match(SET_TX_META);
				setState(180);
				match(LPAREN);
				setState(181);
				((SetTxMetaContext)_localctx).key = match(STRING);
				setState(182);
				match(T__3);
				setState(183);
				((SetTxMetaContext)_localctx).value = expression(0);
				setState(184);
				match(RPAREN);
				}
				break;
			case 4:
				_localctx = new SetAccountMetaContext(_localctx);
				enterOuterAlt(_localctx, 4);
				{
				setState(186);
				match(SET_ACCOUNT_META);
				setState(187);
				match(LPAREN);
				setState(188);
				((SetAccountMetaContext)_localctx).acc = expression(0);
				setState(189);
				match(T__3);
				setState(190);
				((SetAccountMetaContext)_localctx).key = match(STRING);
				setState(191);
				match(T__3);
				setState(192);
				((SetAccountMetaContext)_localctx).value = expression(0);
				setState(193);
				match(RPAREN);
				}
				break;
			case 5:
				_localctx = new FailContext(_localctx);
				enterOuterAlt(_localctx, 5);
				{
				setState(195);
				match(FAIL);
				}
				break;
			case 6:
				_localctx = new SendContext(_localctx);
				enterOuterAlt(_localctx, 6);
				{
				setState(196);
				match(SEND);
				setState(197);
				((SendContext)_localctx).mon = expression(0);
				setState(198);
				match(LPAREN);
				setState(199);
				match(NEWLINE);
				setState(216);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case SOURCE:
					{
					setState(200);
					match(SOURCE);
					setState(201);
					match(EQ);
					setState(202);
					((SendContext)_localctx).src = valueAwareSource();
					setState(203);
					match(NEWLINE);
					setState(204);
					match(DESTINATION);
					setState(205);
					match(EQ);
					setState(206);
					((SendContext)_localctx).dest = destination();
					}
					break;
				case DESTINATION:
					{
					setState(208);
					match(DESTINATION);
					setState(209);
					match(EQ);
					setState(210);
					((SendContext)_localctx).dest = destination();
					setState(211);
					match(NEWLINE);
					setState(212);
					match(SOURCE);
					setState(213);
					match(EQ);
					setState(214);
					((SendContext)_localctx).src = valueAwareSource();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(218);
				match(NEWLINE);
				setState(219);
				match(RPAREN);
				}
				break;
			case 7:
				_localctx = new SendAllContext(_localctx);
				enterOuterAlt(_localctx, 7);
				{
				setState(221);
				match(SEND);
				setState(222);
				((SendAllContext)_localctx).monAll = monetaryAll();
				setState(223);
				match(LPAREN);
				setState(224);
				match(NEWLINE);
				setState(241);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case SOURCE:
					{
					setState(225);
					match(SOURCE);
					setState(226);
					match(EQ);
					setState(227);
					((SendAllContext)_localctx).src = source();
					setState(228);
					match(NEWLINE);
					setState(229);
					match(DESTINATION);
					setState(230);
					match(EQ);
					setState(231);
					((SendAllContext)_localctx).dest = destination();
					}
					break;
				case DESTINATION:
					{
					setState(233);
					match(DESTINATION);
					setState(234);
					match(EQ);
					setState(235);
					((SendAllContext)_localctx).dest = destination();
					setState(236);
					match(NEWLINE);
					setState(237);
					match(SOURCE);
					setState(238);
					match(EQ);
					setState(239);
					((SendAllContext)_localctx).src = source();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(243);
				match(NEWLINE);
				setState(244);
				match(RPAREN);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class Type_Context extends ParserRuleContext {
		public TerminalNode TY_ACCOUNT() { return getToken(NumScriptParser.TY_ACCOUNT, 0); }
		public TerminalNode TY_ASSET() { return getToken(NumScriptParser.TY_ASSET, 0); }
		public TerminalNode TY_NUMBER() { return getToken(NumScriptParser.TY_NUMBER, 0); }
		public TerminalNode TY_STRING() { return getToken(NumScriptParser.TY_STRING, 0); }
		public TerminalNode TY_MONETARY() { return getToken(NumScriptParser.TY_MONETARY, 0); }
		public TerminalNode TY_PORTION() { return getToken(NumScriptParser.TY_PORTION, 0); }
		public Type_Context(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_type_; }
	}

	public final Type_Context type_() throws RecognitionException {
		Type_Context _localctx = new Type_Context(_ctx, getState());
		enterRule(_localctx, 36, RULE_type_);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(248);
			_la = _input.LA(1);
			if ( !((((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << TY_ACCOUNT) | (1L << TY_ASSET) | (1L << TY_NUMBER) | (1L << TY_MONETARY) | (1L << TY_PORTION) | (1L << TY_STRING))) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class OriginContext extends ParserRuleContext {
		public OriginContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_origin; }
	 
		public OriginContext() { }
		public void copyFrom(OriginContext ctx) {
			super.copyFrom(ctx);
		}
	}
	public static class OriginAccountBalanceContext extends OriginContext {
		public ExpressionContext account;
		public ExpressionContext asset;
		public TerminalNode BALANCE() { return getToken(NumScriptParser.BALANCE, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public OriginAccountBalanceContext(OriginContext ctx) { copyFrom(ctx); }
	}
	public static class OriginAccountMetaContext extends OriginContext {
		public ExpressionContext account;
		public Token key;
		public TerminalNode META() { return getToken(NumScriptParser.META, 0); }
		public TerminalNode LPAREN() { return getToken(NumScriptParser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(NumScriptParser.RPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode STRING() { return getToken(NumScriptParser.STRING, 0); }
		public OriginAccountMetaContext(OriginContext ctx) { copyFrom(ctx); }
	}

	public final OriginContext origin() throws RecognitionException {
		OriginContext _localctx = new OriginContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_origin);
		try {
			setState(264);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case META:
				_localctx = new OriginAccountMetaContext(_localctx);
				enterOuterAlt(_localctx, 1);
				{
				setState(250);
				match(META);
				setState(251);
				match(LPAREN);
				setState(252);
				((OriginAccountMetaContext)_localctx).account = expression(0);
				setState(253);
				match(T__3);
				setState(254);
				((OriginAccountMetaContext)_localctx).key = match(STRING);
				setState(255);
				match(RPAREN);
				}
				break;
			case BALANCE:
				_localctx = new OriginAccountBalanceContext(_localctx);
				enterOuterAlt(_localctx, 2);
				{
				setState(257);
				match(BALANCE);
				setState(258);
				match(LPAREN);
				setState(259);
				((OriginAccountBalanceContext)_localctx).account = expression(0);
				setState(260);
				match(T__3);
				setState(261);
				((OriginAccountBalanceContext)_localctx).asset = expression(0);
				setState(262);
				match(RPAREN);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class VarDeclContext extends ParserRuleContext {
		public Type_Context ty;
		public VariableContext name;
		public OriginContext orig;
		public Type_Context type_() {
			return getRuleContext(Type_Context.class,0);
		}
		public VariableContext variable() {
			return getRuleContext(VariableContext.class,0);
		}
		public TerminalNode EQ() { return getToken(NumScriptParser.EQ, 0); }
		public OriginContext origin() {
			return getRuleContext(OriginContext.class,0);
		}
		public VarDeclContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_varDecl; }
	}

	public final VarDeclContext varDecl() throws RecognitionException {
		VarDeclContext _localctx = new VarDeclContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_varDecl);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(266);
			((VarDeclContext)_localctx).ty = type_();
			setState(267);
			((VarDeclContext)_localctx).name = variable();
			setState(270);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==EQ) {
				{
				setState(268);
				match(EQ);
				setState(269);
				((VarDeclContext)_localctx).orig = origin();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class VarListDeclContext extends ParserRuleContext {
		public VarDeclContext varDecl;
		public List<VarDeclContext> v = new ArrayList<VarDeclContext>();
		public TerminalNode VARS() { return getToken(NumScriptParser.VARS, 0); }
		public TerminalNode LBRACE() { return getToken(NumScriptParser.LBRACE, 0); }
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public TerminalNode RBRACE() { return getToken(NumScriptParser.RBRACE, 0); }
		public List<VarDeclContext> varDecl() {
			return getRuleContexts(VarDeclContext.class);
		}
		public VarDeclContext varDecl(int i) {
			return getRuleContext(VarDeclContext.class,i);
		}
		public VarListDeclContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_varListDecl; }
	}

	public final VarListDeclContext varListDecl() throws RecognitionException {
		VarListDeclContext _localctx = new VarListDeclContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_varListDecl);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(272);
			match(VARS);
			setState(273);
			match(LBRACE);
			setState(274);
			match(NEWLINE);
			setState(281); 
			_errHandler.sync(this);
			_la = _input.LA(1);
			do {
				{
				{
				setState(275);
				((VarListDeclContext)_localctx).varDecl = varDecl();
				((VarListDeclContext)_localctx).v.add(((VarListDeclContext)_localctx).varDecl);
				setState(277); 
				_errHandler.sync(this);
				_la = _input.LA(1);
				do {
					{
					{
					setState(276);
					match(NEWLINE);
					}
					}
					setState(279); 
					_errHandler.sync(this);
					_la = _input.LA(1);
				} while ( _la==NEWLINE );
				}
				}
				setState(283); 
				_errHandler.sync(this);
				_la = _input.LA(1);
			} while ( (((_la) & ~0x3f) == 0 && ((1L << _la) & ((1L << TY_ACCOUNT) | (1L << TY_ASSET) | (1L << TY_NUMBER) | (1L << TY_MONETARY) | (1L << TY_PORTION) | (1L << TY_STRING))) != 0) );
			setState(285);
			match(RBRACE);
			setState(286);
			match(NEWLINE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public static class ScriptContext extends ParserRuleContext {
		public VarListDeclContext vars;
		public StatementContext statement;
		public List<StatementContext> stmts = new ArrayList<StatementContext>();
		public TerminalNode EOF() { return getToken(NumScriptParser.EOF, 0); }
		public List<StatementContext> statement() {
			return getRuleContexts(StatementContext.class);
		}
		public StatementContext statement(int i) {
			return getRuleContext(StatementContext.class,i);
		}
		public List<TerminalNode> NEWLINE() { return getTokens(NumScriptParser.NEWLINE); }
		public TerminalNode NEWLINE(int i) {
			return getToken(NumScriptParser.NEWLINE, i);
		}
		public VarListDeclContext varListDecl() {
			return getRuleContext(VarListDeclContext.class,0);
		}
		public ScriptContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_script; }
	}

	public final ScriptContext script() throws RecognitionException {
		ScriptContext _localctx = new ScriptContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_script);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(291);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==NEWLINE) {
				{
				{
				setState(288);
				match(NEWLINE);
				}
				}
				setState(293);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(295);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==VARS) {
				{
				setState(294);
				((ScriptContext)_localctx).vars = varListDecl();
				}
			}

			setState(297);
			((ScriptContext)_localctx).statement = statement();
			((ScriptContext)_localctx).stmts.add(((ScriptContext)_localctx).statement);
			setState(302);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,24,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(298);
					match(NEWLINE);
					setState(299);
					((ScriptContext)_localctx).statement = statement();
					((ScriptContext)_localctx).stmts.add(((ScriptContext)_localctx).statement);
					}
					} 
				}
				setState(304);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,24,_ctx);
			}
			setState(308);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==NEWLINE) {
				{
				{
				setState(305);
				match(NEWLINE);
				}
				}
				setState(310);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(311);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public boolean sempred(RuleContext _localctx, int ruleIndex, int predIndex) {
		switch (ruleIndex) {
		case 4:
			return expression_sempred((ExpressionContext)_localctx, predIndex);
		}
		return true;
	}
	private boolean expression_sempred(ExpressionContext _localctx, int predIndex) {
		switch (predIndex) {
		case 0:
			return precpred(_ctx, 4);
		}
		return true;
	}

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\3\61\u013c\4\2\t\2"+
		"\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4\13"+
		"\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21\t\21\4\22\t\22"+
		"\4\23\t\23\4\24\t\24\4\25\t\25\4\26\t\26\4\27\t\27\4\30\t\30\3\2\3\2\3"+
		"\2\3\2\3\2\3\3\3\3\3\3\3\3\3\3\3\4\3\4\3\4\3\4\3\4\5\4@\n\4\3\5\3\5\3"+
		"\6\3\6\3\6\3\6\5\6H\n\6\3\6\3\6\3\6\7\6M\n\6\f\6\16\6P\13\6\3\7\3\7\3"+
		"\7\5\7U\n\7\3\b\3\b\3\b\3\b\3\b\3\b\3\b\6\b^\n\b\r\b\16\b_\3\b\3\b\3\b"+
		"\3\b\3\b\3\t\3\t\3\t\3\t\3\t\3\t\6\tm\n\t\r\t\16\tn\3\t\3\t\3\n\3\n\3"+
		"\n\5\nv\n\n\3\13\3\13\3\13\5\13{\n\13\3\f\3\f\3\f\5\f\u0080\n\f\3\r\3"+
		"\r\5\r\u0084\n\r\3\16\3\16\3\16\3\16\3\16\6\16\u008b\n\16\r\16\16\16\u008c"+
		"\3\16\3\16\3\17\3\17\3\17\3\17\3\17\3\20\3\20\3\20\5\20\u0099\n\20\3\21"+
		"\3\21\3\21\3\21\3\21\3\21\3\21\6\21\u00a2\n\21\r\21\16\21\u00a3\3\21\3"+
		"\21\3\22\3\22\5\22\u00aa\n\22\3\23\3\23\3\23\3\23\3\23\5\23\u00b1\n\23"+
		"\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23"+
		"\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23"+
		"\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\5\23\u00db"+
		"\n\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23"+
		"\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\3\23\5\23\u00f4\n\23\3\23"+
		"\3\23\3\23\5\23\u00f9\n\23\3\24\3\24\3\25\3\25\3\25\3\25\3\25\3\25\3\25"+
		"\3\25\3\25\3\25\3\25\3\25\3\25\3\25\5\25\u010b\n\25\3\26\3\26\3\26\3\26"+
		"\5\26\u0111\n\26\3\27\3\27\3\27\3\27\3\27\6\27\u0118\n\27\r\27\16\27\u0119"+
		"\6\27\u011c\n\27\r\27\16\27\u011d\3\27\3\27\3\27\3\30\7\30\u0124\n\30"+
		"\f\30\16\30\u0127\13\30\3\30\5\30\u012a\n\30\3\30\3\30\3\30\7\30\u012f"+
		"\n\30\f\30\16\30\u0132\13\30\3\30\7\30\u0135\n\30\f\30\16\30\u0138\13"+
		"\30\3\30\3\30\3\30\2\3\n\31\2\4\6\b\n\f\16\20\22\24\26\30\32\34\36 \""+
		"$&(*,.\2\4\3\2\30\31\3\2!&\2\u014a\2\60\3\2\2\2\4\65\3\2\2\2\6?\3\2\2"+
		"\2\bA\3\2\2\2\nG\3\2\2\2\fT\3\2\2\2\16V\3\2\2\2\20f\3\2\2\2\22u\3\2\2"+
		"\2\24z\3\2\2\2\26\177\3\2\2\2\30\u0081\3\2\2\2\32\u0085\3\2\2\2\34\u0090"+
		"\3\2\2\2\36\u0098\3\2\2\2 \u009a\3\2\2\2\"\u00a9\3\2\2\2$\u00f8\3\2\2"+
		"\2&\u00fa\3\2\2\2(\u010a\3\2\2\2*\u010c\3\2\2\2,\u0112\3\2\2\2.\u0125"+
		"\3\2\2\2\60\61\7\34\2\2\61\62\5\n\6\2\62\63\5\n\6\2\63\64\7\35\2\2\64"+
		"\3\3\2\2\2\65\66\7\34\2\2\66\67\5\n\6\2\678\7\3\2\289\7\35\2\29\5\3\2"+
		"\2\2:@\7\60\2\2;@\7\61\2\2<@\7-\2\2=@\7\'\2\2>@\7(\2\2?:\3\2\2\2?;\3\2"+
		"\2\2?<\3\2\2\2?=\3\2\2\2?>\3\2\2\2@\7\3\2\2\2AB\7/\2\2B\t\3\2\2\2CD\b"+
		"\6\1\2DH\5\6\4\2EH\5\b\5\2FH\5\2\2\2GC\3\2\2\2GE\3\2\2\2GF\3\2\2\2HN\3"+
		"\2\2\2IJ\f\6\2\2JK\t\2\2\2KM\5\n\6\7LI\3\2\2\2MP\3\2\2\2NL\3\2\2\2NO\3"+
		"\2\2\2O\13\3\2\2\2PN\3\2\2\2QU\7(\2\2RU\5\b\5\2SU\7)\2\2TQ\3\2\2\2TR\3"+
		"\2\2\2TS\3\2\2\2U\r\3\2\2\2VW\7\36\2\2W]\7\7\2\2XY\7\24\2\2YZ\5\n\6\2"+
		"Z[\5\22\n\2[\\\7\7\2\2\\^\3\2\2\2]X\3\2\2\2^_\3\2\2\2_]\3\2\2\2_`\3\2"+
		"\2\2`a\3\2\2\2ab\7)\2\2bc\5\22\n\2cd\7\7\2\2de\7\37\2\2e\17\3\2\2\2fg"+
		"\7\36\2\2gl\7\7\2\2hi\5\f\7\2ij\5\22\n\2jk\7\7\2\2km\3\2\2\2lh\3\2\2\2"+
		"mn\3\2\2\2nl\3\2\2\2no\3\2\2\2op\3\2\2\2pq\7\37\2\2q\21\3\2\2\2rs\7\26"+
		"\2\2sv\5\24\13\2tv\7*\2\2ur\3\2\2\2ut\3\2\2\2v\23\3\2\2\2w{\5\n\6\2x{"+
		"\5\16\b\2y{\5\20\t\2zw\3\2\2\2zx\3\2\2\2zy\3\2\2\2{\25\3\2\2\2|}\7\4\2"+
		"\2}\u0080\5\n\6\2~\u0080\7\5\2\2\177|\3\2\2\2\177~\3\2\2\2\u0080\27\3"+
		"\2\2\2\u0081\u0083\5\n\6\2\u0082\u0084\5\26\f\2\u0083\u0082\3\2\2\2\u0083"+
		"\u0084\3\2\2\2\u0084\31\3\2\2\2\u0085\u0086\7\36\2\2\u0086\u008a\7\7\2"+
		"\2\u0087\u0088\5\36\20\2\u0088\u0089\7\7\2\2\u0089\u008b\3\2\2\2\u008a"+
		"\u0087\3\2\2\2\u008b\u008c\3\2\2\2\u008c\u008a\3\2\2\2\u008c\u008d\3\2"+
		"\2\2\u008d\u008e\3\2\2\2\u008e\u008f\7\37\2\2\u008f\33\3\2\2\2\u0090\u0091"+
		"\7\24\2\2\u0091\u0092\5\n\6\2\u0092\u0093\7\23\2\2\u0093\u0094\5\36\20"+
		"\2\u0094\35\3\2\2\2\u0095\u0099\5\30\r\2\u0096\u0099\5\34\17\2\u0097\u0099"+
		"\5\32\16\2\u0098\u0095\3\2\2\2\u0098\u0096\3\2\2\2\u0098\u0097\3\2\2\2"+
		"\u0099\37\3\2\2\2\u009a\u009b\7\36\2\2\u009b\u00a1\7\7\2\2\u009c\u009d"+
		"\5\f\7\2\u009d\u009e\7\23\2\2\u009e\u009f\5\36\20\2\u009f\u00a0\7\7\2"+
		"\2\u00a0\u00a2\3\2\2\2\u00a1\u009c\3\2\2\2\u00a2\u00a3\3\2\2\2\u00a3\u00a1"+
		"\3\2\2\2\u00a3\u00a4\3\2\2\2\u00a4\u00a5\3\2\2\2\u00a5\u00a6\7\37\2\2"+
		"\u00a6!\3\2\2\2\u00a7\u00aa\5\36\20\2\u00a8\u00aa\5 \21\2\u00a9\u00a7"+
		"\3\2\2\2\u00a9\u00a8\3\2\2\2\u00aa#\3\2\2\2\u00ab\u00ac\7\17\2\2\u00ac"+
		"\u00f9\5\n\6\2\u00ad\u00b0\7,\2\2\u00ae\u00b1\5\n\6\2\u00af\u00b1\5\4"+
		"\3\2\u00b0\u00ae\3\2\2\2\u00b0\u00af\3\2\2\2\u00b1\u00b2\3\2\2\2\u00b2"+
		"\u00b3\7\23\2\2\u00b3\u00b4\5\n\6\2\u00b4\u00f9\3\2\2\2\u00b5\u00b6\7"+
		"\r\2\2\u00b6\u00b7\7\32\2\2\u00b7\u00b8\7\'\2\2\u00b8\u00b9\7\6\2\2\u00b9"+
		"\u00ba\5\n\6\2\u00ba\u00bb\7\33\2\2\u00bb\u00f9\3\2\2\2\u00bc\u00bd\7"+
		"\16\2\2\u00bd\u00be\7\32\2\2\u00be\u00bf\5\n\6\2\u00bf\u00c0\7\6\2\2\u00c0"+
		"\u00c1\7\'\2\2\u00c1\u00c2\7\6\2\2\u00c2\u00c3\5\n\6\2\u00c3\u00c4\7\33"+
		"\2\2\u00c4\u00f9\3\2\2\2\u00c5\u00f9\7\20\2\2\u00c6\u00c7\7\21\2\2\u00c7"+
		"\u00c8\5\n\6\2\u00c8\u00c9\7\32\2\2\u00c9\u00da\7\7\2\2\u00ca\u00cb\7"+
		"\22\2\2\u00cb\u00cc\7 \2\2\u00cc\u00cd\5\"\22\2\u00cd\u00ce\7\7\2\2\u00ce"+
		"\u00cf\7\25\2\2\u00cf\u00d0\7 \2\2\u00d0\u00d1\5\24\13\2\u00d1\u00db\3"+
		"\2\2\2\u00d2\u00d3\7\25\2\2\u00d3\u00d4\7 \2\2\u00d4\u00d5\5\24\13\2\u00d5"+
		"\u00d6\7\7\2\2\u00d6\u00d7\7\22\2\2\u00d7\u00d8\7 \2\2\u00d8\u00d9\5\""+
		"\22\2\u00d9\u00db\3\2\2\2\u00da\u00ca\3\2\2\2\u00da\u00d2\3\2\2\2\u00db"+
		"\u00dc\3\2\2\2\u00dc\u00dd\7\7\2\2\u00dd\u00de\7\33\2\2\u00de\u00f9\3"+
		"\2\2\2\u00df\u00e0\7\21\2\2\u00e0\u00e1\5\4\3\2\u00e1\u00e2\7\32\2\2\u00e2"+
		"\u00f3\7\7\2\2\u00e3\u00e4\7\22\2\2\u00e4\u00e5\7 \2\2\u00e5\u00e6\5\36"+
		"\20\2\u00e6\u00e7\7\7\2\2\u00e7\u00e8\7\25\2\2\u00e8\u00e9\7 \2\2\u00e9"+
		"\u00ea\5\24\13\2\u00ea\u00f4\3\2\2\2\u00eb\u00ec\7\25\2\2\u00ec\u00ed"+
		"\7 \2\2\u00ed\u00ee\5\24\13\2\u00ee\u00ef\7\7\2\2\u00ef\u00f0\7\22\2\2"+
		"\u00f0\u00f1\7 \2\2\u00f1\u00f2\5\36\20\2\u00f2\u00f4\3\2\2\2\u00f3\u00e3"+
		"\3\2\2\2\u00f3\u00eb\3\2\2\2\u00f4\u00f5\3\2\2\2\u00f5\u00f6\7\7\2\2\u00f6"+
		"\u00f7\7\33\2\2\u00f7\u00f9\3\2\2\2\u00f8\u00ab\3\2\2\2\u00f8\u00ad\3"+
		"\2\2\2\u00f8\u00b5\3\2\2\2\u00f8\u00bc\3\2\2\2\u00f8\u00c5\3\2\2\2\u00f8"+
		"\u00c6\3\2\2\2\u00f8\u00df\3\2\2\2\u00f9%\3\2\2\2\u00fa\u00fb\t\3\2\2"+
		"\u00fb\'\3\2\2\2\u00fc\u00fd\7\f\2\2\u00fd\u00fe\7\32\2\2\u00fe\u00ff"+
		"\5\n\6\2\u00ff\u0100\7\6\2\2\u0100\u0101\7\'\2\2\u0101\u0102\7\33\2\2"+
		"\u0102\u010b\3\2\2\2\u0103\u0104\7+\2\2\u0104\u0105\7\32\2\2\u0105\u0106"+
		"\5\n\6\2\u0106\u0107\7\6\2\2\u0107\u0108\5\n\6\2\u0108\u0109\7\33\2\2"+
		"\u0109\u010b\3\2\2\2\u010a\u00fc\3\2\2\2\u010a\u0103\3\2\2\2\u010b)\3"+
		"\2\2\2\u010c\u010d\5&\24\2\u010d\u0110\5\b\5\2\u010e\u010f\7 \2\2\u010f"+
		"\u0111\5(\25\2\u0110\u010e\3\2\2\2\u0110\u0111\3\2\2\2\u0111+\3\2\2\2"+
		"\u0112\u0113\7\13\2\2\u0113\u0114\7\36\2\2\u0114\u011b\7\7\2\2\u0115\u0117"+
		"\5*\26\2\u0116\u0118\7\7\2\2\u0117\u0116\3\2\2\2\u0118\u0119\3\2\2\2\u0119"+
		"\u0117\3\2\2\2\u0119\u011a\3\2\2\2\u011a\u011c\3\2\2\2\u011b\u0115\3\2"+
		"\2\2\u011c\u011d\3\2\2\2\u011d\u011b\3\2\2\2\u011d\u011e\3\2\2\2\u011e"+
		"\u011f\3\2\2\2\u011f\u0120\7\37\2\2\u0120\u0121\7\7\2\2\u0121-\3\2\2\2"+
		"\u0122\u0124\7\7\2\2\u0123\u0122\3\2\2\2\u0124\u0127\3\2\2\2\u0125\u0123"+
		"\3\2\2\2\u0125\u0126\3\2\2\2\u0126\u0129\3\2\2\2\u0127\u0125\3\2\2\2\u0128"+
		"\u012a\5,\27\2\u0129\u0128\3\2\2\2\u0129\u012a\3\2\2\2\u012a\u012b\3\2"+
		"\2\2\u012b\u0130\5$\23\2\u012c\u012d\7\7\2\2\u012d\u012f\5$\23\2\u012e"+
		"\u012c\3\2\2\2\u012f\u0132\3\2\2\2\u0130\u012e\3\2\2\2\u0130\u0131\3\2"+
		"\2\2\u0131\u0136\3\2\2\2\u0132\u0130\3\2\2\2\u0133\u0135\7\7\2\2\u0134"+
		"\u0133\3\2\2\2\u0135\u0138\3\2\2\2\u0136\u0134\3\2\2\2\u0136\u0137\3\2"+
		"\2\2\u0137\u0139\3\2\2\2\u0138\u0136\3\2\2\2\u0139\u013a\7\2\2\3\u013a"+
		"/\3\2\2\2\34?GNT_nuz\177\u0083\u008c\u0098\u00a3\u00a9\u00b0\u00da\u00f3"+
		"\u00f8\u010a\u0110\u0119\u011d\u0125\u0129\u0130\u0136";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}