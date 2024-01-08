<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class WorkflowInstanceHistoryStageOutput
{
	#[\JMS\Serializer\Annotation\SerializedName('CreateTransaction')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityCreateTransactionOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityCreateTransactionOutput $createTransaction = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('DebitWallet')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityDebitWalletOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityDebitWalletOutput $debitWallet = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('GetAccount')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityGetAccountOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityGetAccountOutput $getAccount = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('GetPayment')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityGetPaymentOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityGetPaymentOutput $getPayment = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('GetWallet')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityGetWalletOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityGetWalletOutput $getWallet = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('ListWallets')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\OrchestrationListWalletsResponse')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?OrchestrationListWalletsResponse $listWallets = null;
    
	#[\JMS\Serializer\Annotation\SerializedName('RevertTransaction')]
    #[\JMS\Serializer\Annotation\Type('formance\stack\Models\Shared\ActivityRevertTransactionOutput')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?ActivityRevertTransactionOutput $revertTransaction = null;
    
	public function __construct()
	{
		$this->createTransaction = null;
		$this->debitWallet = null;
		$this->getAccount = null;
		$this->getPayment = null;
		$this->getWallet = null;
		$this->listWallets = null;
		$this->revertTransaction = null;
	}
}
